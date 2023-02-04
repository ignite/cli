# Gönderileri listele

Bu bölümde, kullanıcıların blok zinciri uygulamanızda depolanan tüm blog gönderilerini almalarını sağlayan bir özellik geliştireceksiniz. Bu özellik, kullanıcıların bir sorgu gerçekleştirmesine ve sayfalandırılmış bir yanıt almasına olanak tanıyacaktır; bu, çıktının daha küçük veri parçalarına veya "sayfalara" bölüneceği anlamına gelir. Bu sayede kullanıcılar, potansiyel olarak uzun bir listeyi tek seferde kaydırmak yerine bir seferde belirli sayıda gönderiyi görüntüleyebilecekleri için gönderi listesinde daha kolay gezinebilecek ve göz atabileceklerdir.

### Gönderileri listele

Bir kullanıcı blok zinciri uygulamasına bir sorgu yaptığında çağrılacak olan `ListPost` keeper yöntemini uygulayalım ve zincirde depolanan tüm gönderilerin sayfalandırılmış bir listesini talep edelim.

x/blog/keeper/query\_list\_post.go

```
package keeper

import (
    "context"

    "blog/x/blog/types"

    "github.com/cosmos/cosmos-sdk/store/prefix"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/query"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (k Keeper) ListPost(goCtx context.Context, req *types.QueryListPostRequest) (*types.QueryListPostResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    var posts []types.Post
    ctx := sdk.UnwrapSDKContext(goCtx)

    store := ctx.KVStore(k.storeKey)
    postStore := prefix.NewStore(store, types.KeyPrefix(types.PostKey))

    pageRes, err := query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
        var post types.Post
        if err := k.cdc.Unmarshal(value, &post); err != nil {
            return err
        }

        posts = append(posts, post)
        return nil
    })

    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &types.QueryListPostResponse{Post: posts, Pagination: pageRes}, nil
}
```

`ListPost` iki bağımsız değişken alır: bir context nesnesi ve `QueryListPostRequest` türünde bir istek nesnesi. `QueryListPostResponse` türünde bir yanıt nesnesi ve bir hata döndürür.

Fonksiyon ilk olarak istek nesnesinin `nil` olup olmadığını kontrol eder ve nil ise `InvalidArgument` koduyla bir hata döndürür. Daha sonra `Post` nesnelerinin boş bir dilimini başlatır ve bağlam nesnesini açar.

Keeper struct'ın `storeKey` alanını kullanarak bağlamdan bir anahtar-değer deposu alır ve `PostKey` önekini kullanarak yeni bir depo oluşturur. Ardından, mağaza ve istek nesnesindeki sayfalama bilgileri üzerinde `query` paketinden `Paginate` işlevini çağırır. Paginate'e argüman olarak aktarılan fonksiyon, depodaki anahtar-değer çiftleri üzerinde yineleme yapar ve değerleri `Post` nesnelerine ayırır; bu nesneler daha sonra `posts` dilimine eklenir.

Sayfalandırma sırasında bir hata oluşursa, fonksiyon hata mesajıyla birlikte bir `Internal error` döndürür. Aksi takdirde, gönderilerin listesini ve sayfalama bilgilerini içeren bir `QueryListPostResponse` nesnesi döndürür.

### `QueryListPostResponse`'u Değiştirme

Gönderilerin bir listesini döndürmek için tekrarlanan bir anahtar kelime ekleyin ve alanı işaretçi olmadan oluşturmak için \[`(gogoproto.nullable) = false`] seçeneğini ekleyin.

proto/blog/blog/query.proto

```
message QueryListPostResponse {
  repeated Post post = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

Proto'dan Go dosyaları oluşturmak için komutu çalıştırın:

```
ignite generate proto-go
```
