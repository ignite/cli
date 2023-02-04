# cosmostxcollector

Paket, Cosmos blok zincirlerinden işlemlerin ve olayların bir veri arka ucuna toplanması için destek uygular ve ayrıca toplanan verilerin sorgulanması için destek ekler.

### İşlem ve olay verisi toplama

İşlemler ve olaylar `cosmostxcollector.Collector` türü kullanılarak toplanabilir. Bu tür, her bloktan veri almak için bir `cosmosclient.Client` örneği ve verileri kaydetmek için bir veri arka uç adaptörü kullanır.

#### Veri arka uç bağdaştırıcıları[​](broken-reference) <a href="#data-backend-adapters" id="data-backend-adapters"></a>

Veri arka uç adaptörleri, toplanan verileri sorgulamak ve farklı veri arka uçlarına kaydetmek için kullanılır ve `cosmostxcollector.adapter.Adapter` arayüzünü uygulamalıdır.

PostgreSQL için bir adaptör `cosmostxcollector.adapter.postgres.Adapter`'de zaten uygulanmıştır. Örneklerde kullanılan budur.

#### Örnek: Veri toplama[​](broken-reference) <a href="#example-data-collection" id="example-data-collection"></a>

Veri toplama örneği, yerel ortamda çalışan ve "cosmos" adında boş bir veritabanı içeren bir PostgreSQL veritabanı olduğunu varsayar.

Gerekli veritabanı tabloları, ilk kez çalıştırıldığında toplayıcı tarafından otomatik olarak oluşturulacaktır.

Uygulama çalıştırıldığında, son bloklardan birinden başlayarak mevcut blok yüksekliğine kadar olan tüm işlemleri ve olayları getirecek ve veritabanını dolduracaktır:

```
package main

import (
    "context"
    "log"

    "github.com/ignite/cli/ignite/pkg/clictx"
    "github.com/ignite/cli/ignite/pkg/cosmosclient"
    "github.com/ignite/cli/ignite/pkg/cosmostxcollector"
    "github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
)

const (
    // Name of a local PostgreSQL database
    dbName = "cosmos"

    // Cosmos RPC address
    rpcAddr = "https://rpc.cosmos.network:443"
)

func collect(ctx context.Context, db postgres.Adapter) error {
    // Make sure that the data backend schema is up to date
    if err = db.Init(ctx); err != nil {
        return err
    }

    // Init the Cosmos client
    client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr))
    if err != nil {
        return err
    }

    // Get the latest block height
    latestHeight, err := client.LatestBlockHeight(ctx)
    if err != nil {
        return err
    }

    // Collect transactions and events starting from a block height.
    // The collector stops at the latest height available at the time of the call.
    collector := cosmostxcollector.New(db, client)
    if err := collector.Collect(ctx, latestHeight-50); err != nil {
        return err
    }

    return nil
}

func main() {
    ctx := clictx.From(context.Background())

    // Init an adapter for a local PostgreSQL database running with the default values
    params := map[string]string{"sslmode": "disable"}
    db, err := postgres.NewAdapter(dbName, postgres.WithParams(params))
    if err != nil {
        log.Fatal(err)
    }

    if err := collect(ctx, db); err != nil {
        log.Fatal(err)
    }
}
```

Toplanan veriler, olay sorguları veya imleç tabanlı sorgular kullanılarak veri arka uç bağdaştırıcıları aracılığıyla sorgulanabilir.

Sorgular, oluşturma sırasında farklı seçenekler kullanarak sıralama, sayfalama ve filtrelemeyi destekler. İmleç tabanlı olanlar ayrıca belirli alanların veya özelliklerin seçilmesini ve sorgunun bir fonksiyon olduğu durumlarda argümanların iletilmesini de destekler.

Varsayılan olarak sorgulara sıralama, filtreleme veya sayfalama uygulanmaz.

#### Olay sorguları[​](broken-reference) <a href="#event-queries" id="event-queries"></a>

Olay sorguları olayları ve özniteliklerini `[]cosmostxcollector.query.Event` olarak döndürür.

#### Örnek: Olayları sorgulama <a href="#example-query-events" id="example-query-events"></a>

Örnek, Cosmos'un banka modülünden transfer olaylarını okur ve sonuçları sayfalandırır.

```
import (
    "context"

    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
    "github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
    "github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

func queryBankTransferEvents(ctx context.Context, db postgres.Adapter) ([]query.Event, error) {
    // Create an event query that returns events of type "transfer"
    qry := query.NewEventQuery(
        query.WithFilters(
            // Filter transfer events from Cosmos' bank module
            postgres.FilterByEventType(banktypes.EventTypeTransfer),
        ),
        query.WithPageSize(10),
        query.AtPage(1),
    )

    // Execute the query
    return db.QueryEvents(ctx, qry)
}
```

#### İmleç tabanlı sorgular <a href="#cursor-based-queries" id="cursor-based-queries"></a>

Bu tür sorgular, Olay sorgularının kullanışlı olmadığı bağlamlarda kullanılmak üzere tasarlanmıştır.

İmleç tabanlı sorgular, ilişkisel veritabanlarında bir tablo, görünüm veya işlev ya da ilişkisel olmayan veri arka uçlarında bir koleksiyon veya işlev olabilen tek bir "varlığı" sorgulayabilir.

Bu tür sorguların sonucu `cosmostxcollector.query.Cursor` arayüzünü uygulayan bir imleçtir.

#### Örnek: İmleçleri kullanarak olayları sorgulama[​](broken-reference) <a href="#example-query-events-using-cursors" id="example-query-events-using-cursors"></a>

```
import (
    "context"

    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
    "github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
    "github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

func queryBankTransferEventIDs(ctx context.Context, db postgres.Adapter) (ids []int64, err error) {
    // Create a query that returns the IDs for events of type "transfer"
    qry := query.New(
        "event",
        query.Fields("id"),
        query.WithFilters(
            // Filter transfer events from Cosmos' bank module
            postgres.NewFilter("type", banktypes.EventTypeTransfer),
        ),
        query.WithPageSize(10),
        query.AtPage(1),
        query.SortByFields(query.SortOrderAsc, "id"),
    )

    // Execute the query
    cr, err := db.Query(ctx, qry)
    if err != nil {
        return nil, err
    }

    // Read the results
    for cr.Next() {
        var eventID int64

        if err := cr.Scan(&eventID); err != nil {
            return nil, err
        }

        ids = append(ids, eventID)
    }

    return ids, nil
}
```
