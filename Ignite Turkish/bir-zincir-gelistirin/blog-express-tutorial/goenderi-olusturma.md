# Gönderi oluşturma

Bu bölümde, bir "create post" mesajını işleme sürecine odaklanacağız. Bu işlem, keeper metodu olarak bilinen özel bir fonksiyon türünün kullanımını içerir. [Keeper](https://docs.cosmos.network/main/building-modules/keeper) metotları, blockchain ile etkileşime girmekten ve bir mesajda verilen talimatlara göre durumunu değiştirmekten sorumludur.

Bir "create post" mesajı alındığında, ilgili keeper metodu çağrılır ve mesaj bir argüman olarak iletilir. Keeper yöntemi daha sonra blockchain'in mevcut durumunu almak ve değiştirmek için store nesnesi tarafından sağlanan çeşitli getter ve setter işlevlerini kullanabilir. Bu, keeper yönteminin "create post" mesajını etkin bir şekilde işlemesini ve blok zincirinde gerekli güncellemeleri yapmasını sağlar.

Store nesnesine erişme ve değiştirme kodunu temiz ve keeper yöntemlerinde uygulanan mantıktan ayrı tutmak için `post.go` adında yeni bir dosya oluşturacağız. Bu dosya, blok zincirinde mesajların oluşturulması ve yönetilmesiyle ilgili işlemleri gerçekleştirmek için özel olarak tasarlanmış fonksiyonları içerecektir.

### Mağazaya gönderi ekleme

x/blog/keeper/post.go

```
package keeper

import (
    "encoding/binary"

    "blog/x/blog/types"

    "github.com/cosmos/cosmos-sdk/store/prefix"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AppendPost(ctx sdk.Context, post types.Post) uint64 {
    count := k.GetPostCount(ctx)
    post.Id = count
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
    appendedValue := k.cdc.MustMarshal(&post)
    store.Set(GetPostIDBytes(post.Id), appendedValue)
    k.SetPostCount(ctx, count+1)
    return count
}
```

Bu kod, bir Keeper türüne ait `AppendPost` adlı bir işlevi tanımlar. `Keeper` tipi, blok zinciri ile etkileşime girmekten ve çeşitli mesajlara yanıt olarak durumunu değiştirmekten sorumludur.

AppendPost işlevi iki argüman alır: bir `Context` nesnesi ve bir `Post` nesnesi. [Context](https://docs.cosmos.network/main/core/context) nesnesi, Cosmos SDK'daki birçok işlevde standart bir parametredir ve mevcut blok yüksekliği gibi blok zincirinin mevcut durumu hakkında bağlamsal bilgi sağlamak için kullanılır. `Post` nesnesi, blok zincirine eklenecek bir gönderiyi temsil eder.

İşlev, `GetPostCount` yöntemini kullanarak mevcut gönderi sayısını alarak başlar. Henüz uygulanmadığı için bu yöntemi bir sonraki adımda uygulayacaksınız. Bu yöntem `Keeper` nesnesi üzerinde çağrılır ve bir `Context` nesnesini argüman olarak alır. Blok zincirine eklenen mevcut gönderi sayısını döndürür.

Ardından, işlev yeni gönderinin kimliğini geçerli gönderi sayısı olarak ayarlar, böylece her gönderinin benzersiz bir tanımlayıcısı olur. Bunu count değerini `Post` nesnesinin `Id` alanına atayarak yapar.

Fonksiyon daha sonra `prefix.NewStore` fonksiyonunu kullanarak yeni bir [store](https://docs.cosmos.network/main/core/store) nesnesi oluşturur. `prefix.NewStore` fonksiyonu iki argüman alır: sağlanan context ile ilişkili `KVStore` ve `Post` nesneleri için bir anahtar öneki. `KVStore`, blok zincirindeki verileri kalıcı hale getirmek için kullanılan bir anahtar-değer deposudur ve anahtar öneki, `Post` nesnelerini aynı `KVStore`'da depolanabilecek diğer nesne türlerinden ayırmak için kullanılır.

Fonksiyon daha sonra `cdc.MustMarshal` fonksiyonunu kullanarak `Post` nesnesini serileştirir ve store nesnesinin `Set` metodunu kullanarak blok zincirinde saklar. `cdc.MustMarshal` işlevi Cosmos SDK'nın [kodlama/kod çözme](https://docs.cosmos.network/main/core/encoding) kütüphanesinin bir parçasıdır ve Post nesnesini `KVStore`'da saklanabilecek bir byte dilimine dönüştürmek için kullanılır. `Set` yöntemi store nesnesi üzerinde çağrılır ve iki argüman alır: bir anahtar ve bir değer. Bu durumda, anahtar `GetPostIDBytes` işlevi tarafından oluşturulan bir byte dilimidir ve değer de serileştirilmiş `Post` nesnesidir. Henüz uygulanmadığı için bu yöntemi bir sonraki adımda uygulayacaksınız.

Son olarak, fonksiyon post sayısını bir artırır ve `SetPostCount` yöntemini kullanarak blok zinciri durumunu günceller. Henüz uygulanmadığı için bu yöntemi bir sonraki adımda uygulayacaksınız. Bu yöntem Keeper nesnesi üzerinde çağrılır ve argüman olarak bir `Context` nesnesi ile yeni bir gönderi sayısı alır. Blok zincirindeki mevcut gönderi sayısını, sağlanan yeni gönderi sayısı olacak şekilde günceller.

Fonksiyon daha sonra, artırılmadan önceki mevcut gönderi sayısı olan yeni oluşturulan gönderinin kimliğini döndürür. Bu, fonksiyonu çağıranın blok zincirine yeni eklenen gönderinin kimliğini bilmesini sağlar.

`AppendPost` uygulamasını tamamlamak için aşağıdaki görevlerin gerçekleştirilmesi gerekir:

* Veritabanından gönderileri depolamak ve almak için kullanılacak `PostKey`'i tanımlayın.
* Veritabanında depolanan mevcut gönderi sayısını alacak olan `GetPostCount`'u uygulayın.
* Bir gönderi kimliğini byte dizisine dönüştürecek olan `GetPostIDBytes` öğesini uygulayın.
* `SetPostCount`, veritabanında depolanan gönderi sayısını güncelleyin.

#### Gönderi anahtarı öneki

`keys.go` dosyasında `PostKey` önekini aşağıdaki gibi tanımlayalım:

x/blog/types/keys.go

```
const (
    PostKey = "Post/value/"
)
```

Bu önek, bir gönderiyi benzersiz bir şekilde tanımlamak için kullanılacaktır. sistemi. Her gönderi için anahtarın başlangıcı olarak kullanılacak ve ardından her gönderi için benzersiz bir anahtar oluşturmak üzere gönderinin kimliği gelecektir.

#### Gönderi sayısını alma

`post.go` dosyasında `GetPostCount` fonksiyonunu aşağıdaki gibi tanımlayalım:

x/blog/keeper/post.go

```
func (k Keeper) GetPostCount(ctx sdk.Context) uint64 {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
    byteKey := types.KeyPrefix(types.PostCountKey)
    bz := store.Get(byteKey)
    if bz == nil {
        return 0
    }
    return binary.BigEndian.Uint64(bz)
}
```

Bu kod, `Keepe`r yapısına ait `GetPostCount` adlı bir fonksiyon tanımlar. Fonksiyon tek bir argüman, `sdk.Context` türünde bir context nesnesi `ctx` alır ve `uint64` türünde bir değer döndürür.

Fonksiyon, context'teki anahtar-değer deposunu ve önek olarak boş bir byte dilimini kullanarak yeni bir depo oluşturarak başlar. Ardından, `types` paketindeki `KeyPrefix` işlevini kullanarak bir byte dilimi `byteKey` tanımlar ve `PostCountKey` değerini alır. `PostCountKey`'i bir sonraki adımda tanımlayacaksınız.

Fonksiyon daha sonra `Get` yöntemini kullanarak depodaki `byteKey` anahtarındaki değeri alır ve `bz` değişkeninde saklar.

Ardından, fonksiyon bir if deyimi kullanarak `byteKey` adresindeki değerin `nil` olup olmadığını kontrol eder. Değer `nil` ise, yani anahtar depoda mevcut değilse, fonksiyon 0 değerini döndürür. Bu, anahtarla ilişkili hiçbir öğe veya gönderi olmadığını gösterir.

`byteKey`'deki değer nil değilse, fonksiyon `bz`'deki byte'ları ayrıştırmak için binary paketinin `BigEndian` türünü kullanır ve sonuçta `uint64` değerini döndürür. `BigEndian` türü, `bz` içindeki byte'ları big-endian kodlanmış işaretsiz 64 bit tamsayı olarak yorumlamak için kullanılır. `Uint64` yöntemi byte'ları bir `uint64` değerine dönüştürür ve döndürür.

`GetPostCount` fonksiyonu, anahtar-değer deposunda saklanan ve `uint64` değeri olarak temsil edilen toplam gönderi sayısını almak için kullanılır.

`keys.go` dosyasında `PostCountKey`'i aşağıdaki gibi tanımlayalım:

x/blog/types/keys.go

```
const (
    PostCountKey = "Post/count/"
)
```

Bu anahtar, mağazaya eklenen en son gönderinin kimliğini takip etmek için kullanılacaktır.

#### Gönderi kimliğini byte'a dönüştürme

Şimdi, bir gönderi kimliğini byte dizisine dönüştürecek olan `GetPostIDBytes`'ı uygulayalım.

x/blog/keeper/post.go

```
func GetPostIDBytes(id uint64) []byte {
    bz := make([]byte, 8)
    binary.BigEndian.PutUint64(bz, id)
    return bz
}
```

`GetPostIDBytes`, `uint64` türünde bir `id` değeri alır ve `[]byte` türünde bir değer döndürür.

İşlev, `make` yerleşik işlevini kullanarak 8 uzunluğunda yeni bir byte dilimi `bz` oluşturarak başlar. Ardından, `id` değerini big-endian kodlanmış işaretsiz tamsayı olarak kodlamak için `binary` paketinin `BigEndian` türünü kullanır ve `PutUint64` yöntemini kullanarak sonucu `bz` içinde saklar. Son olarak, fonksiyon elde edilen byte dilimi `bz`'yi döndürür.

Bu fonksiyon, `uint64` olarak temsil edilen bir posta kimliğini, bir anahtar-değer deposunda anahtar olarak kullanılabilecek bir byte dilimine dönüştürmek için kullanılabilir. `binary.BigEndian.PutUint64` fonksiyonu, `id`'nin `uint64` değerini big-endian kodlanmış işaretsiz tamsayı olarak kodlar ve elde edilen byte'ları `[]byte` dilimi `bz`'de saklar. Elde edilen byte dilimi daha sonra depoda bir anahtar olarak kullanılabilir.

#### Gönderi sayısını güncelleme

Veritabanında depolanan gönderi sayısını güncelleyecek olan `post.go`'da `SetPostCount`'u uygulayın.

x/blog/keeper/post.go

```
func (k Keeper) SetPostCount(ctx sdk.Context, count uint64) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
    byteKey := types.KeyPrefix(types.PostCountKey)
    bz := make([]byte, 8)
    binary.BigEndian.PutUint64(bz, count)
    store.Set(byteKey, bz)
}
```

Bu kod, `Keeper` struct içinde `SetPostCount` fonksiyonunu tanımlar. Fonksiyon, `sdk.Context` türünde bir context `ctx` ve `uint64` türünde bir `count` değeri alır ve bir değer döndürmez.

Fonksiyon ilk olarak prefix paketinden `NewStore` fonksiyonunu çağırarak ve context'ten anahtar-değer deposunu ve önek olarak boş bir byte dilimini aktararak yeni bir depo oluşturur. Ortaya çıkan depoyu store adlı bir değişkende saklar.

Ardından, işlev `types` paketindeki `KeyPrefix` işlevini kullanarak ve `PostCountKey` değerini aktararak bir byte dilimi `byteKey` tanımlar. `KeyPrefix` fonksiyonu, önek olarak verilen anahtarla birlikte bir byte dilimi döndürür.

Fonksiyon daha sonra `make` yerleşik fonksiyonunu kullanarak 8 uzunluğunda yeni bir byte dilimi bz oluşturur. Ardından, `count` değerini big-endian kodlanmış işaretsiz tamsayı olarak kodlamak için `binary` paketinin `BigEndian` türünü kullanır ve `PutUint64` yöntemini kullanarak sonucu `bz` içinde saklar.

Son olarak, fonksiyon `store` değişkeni üzerinde `Set` yöntemini çağırır ve `byteKey` ile `bz` değerlerini argüman olarak iletir. Bu, depodaki `byteKey` anahtarındaki değeri bz değerine ayarlar.

Bu fonksiyon, veritabanında depolanan gönderilerin sayısını güncellemek için kullanılabilir. Bunu, `binary.BigEndian.PutUint64` işlevini kullanarak `count` öğesinin `uint64` değerini bir byte dilimine dönüştürerek ve ardından elde edilen byte dilimini `Set` yöntemini kullanarak depodaki `byteKey` anahtarında saklayarak yapar.

Artık blog gönderileri oluşturmak için kodu uyguladığınıza göre, "create post" mesajı işlendiğinde çağrılan keeper yöntemini uygulamaya devam edebilirsiniz.

### "Gönderi oluştur" mesajının işlenmesi

x/blog/keeper/msg\_server\_create\_post.go

```
package keeper

import (
    "context"

    "blog/x/blog/types"

    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    var post = types.Post{
        Creator: msg.Creator,
        Title:   msg.Title,
        Body:    msg.Body,
    }
    id := k.AppendPost(
        ctx,
        post,
    )
    return &types.MsgCreatePostResponse{
        Id: id,
    }, nil
}
```

`CreatePost` işlevi, `MsgCreatePost` mesaj türü için bir mesaj işleyicisidir. `MsgCreatePost` mesajında sağlanan bilgilere dayanarak blok zincirinde yeni bir gönderi oluşturmaktan sorumludur.

Fonksiyon ilk olarak `sdk.UnwrapSDKContext` fonksiyonunu kullanarak Go bağlamından Cosmos SDK bağlamını alır. Ardından `MsgCreatePost` mesajındaki `Creator`, `Title` ve `Body` alanlarını kullanarak yeni bir `Post` nesnesi oluşturur.

Ardından, fonksiyon `msgServer` nesnesinde (Keeper türünde olan) `AppendPost` yöntemini çağırır ve Cosmos SDK context'ini ve yeni `Post` nesnesini argüman olarak iletir. `AppendPost` yöntemi, yeni gönderinin blok zincirine eklenmesinden ve yeni gönderinin kimliğinin döndürülmesinden sorumludur.

Son olarak, fonksiyon yeni gönderinin kimliğini içeren bir `MsgCreatePostResponse` nesnesi döndürür. Ayrıca, işlemin başarılı olduğunu gösteren bir nil hatası döndürür.

### Özet

Harika bir iş çıkardınız! Blog gönderilerini blok zinciri deposuna yazma mantığını ve bir "gönderi oluştur" mesajı işlendiğinde çağrılacak keeper yöntemini başarıyla uyguladınız.

`AppendPost` keeper yöntemi geçerli gönderi sayısını alır, yeni gönderinin ID'sini geçerli gönderi sayısı olarak ayarlar, gönderi nesnesini serileştirir ve `store` nesnesinin `Set` yöntemini kullanarak blok zincirinde saklar. Depodaki gönderinin anahtarı `GetPostIDBytes` işlevi tarafından oluşturulan bir byte dilimidir ve değeri de serileştirilmiş gönderi nesnesidir. Fonksiyon daha sonra gönderi sayısını bir artırır ve `SetPostCount` yöntemini kullanarak blok zinciri durumunu günceller.

`CreatePost` işleyici yöntemi, yeni gönderi için verileri içeren bir `MsgCreatePost` mesajı alır, bu verileri kullanarak yeni bir `Post` nesnesi oluşturur ve blok zincirine eklenmek üzere `AppendPost` keeper yöntemine iletir. Ardından, yeni oluşturulan gönderinin kimliğini içeren bir `MsgCreatePostResponse` nesnesi döndürür.

Bu yöntemleri uygulayarak, "create post" mesajlarını işlemek ve blok zincirine gönderi eklemek için gerekli mantığı başarıyla uygulamış olursunuz.
