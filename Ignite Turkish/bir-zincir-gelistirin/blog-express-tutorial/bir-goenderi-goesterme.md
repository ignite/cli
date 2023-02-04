# Bir gönderi gösterme

Bu bölümde, blog uygulamanızda kullanıcıların tek tek blog gönderilerini benzersiz kimliklerine göre almalarını sağlayan bir özellik uygulayacaksınız. Bu kimlik, her blog gönderisine oluşturulduğunda atanır ve blok zincirinde saklanır. Bu sorgulama işlevini ekleyerek, kullanıcılar kimliklerini belirterek belirli blog gönderilerini kolayca alabilecekler.

### Gönderiyi göster

Bir kullanıcı blok zinciri uygulamasına bir sorgu yaptığında çağrılacak olan `ShowPost` keeper yöntemini, istenen gönderinin kimliğini belirterek uygulayalım.

x/blog/keeper/query\_show\_post.go

```
package keeper

import (
    "context"

    "blog/x/blog/types"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (k Keeper) ShowPost(goCtx context.Context, req *types.QueryShowPostRequest) (*types.QueryShowPostResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    ctx := sdk.UnwrapSDKContext(goCtx)
    post, found := k.GetPost(ctx, req.Id)
    if !found {
        return nil, sdkerrors.ErrKeyNotFound
    }

    return &types.QueryShowPostResponse{Post: post}, nil
}
```

`ShowPost`, blok zincirinin durumundan tek bir post nesnesi almak için kullanılan bir fonksiyondur. İki argüman alır: `goCtx` adında bir `context.Context` nesnesi ve `req` adında bir `types.QueryShowPostRequest` nesnesine bir işaretçi. Bir `types.QueryShowPostResponse` nesnesine bir işaretçi ve bir `error` döndürür.

Fonksiyon önce `req` bağımsız değişkeninin `nil` olup olmadığını kontrol eder. Eğer öyleyse, `google.golang.org/grpc/status` paketindeki `status.Error` fonksiyonunu kullanarak `InvalidArgument` koduyla ve "invalid request" mesajıyla bir hata döndürür.

Eğer `req` argümanı `nil` değilse, fonksiyon `sdk.UnwrapSDKContext` fonksiyonunu kullanarak `sdk.Context` nesnesini `context.Context` nesnesinden çözer. Daha sonra `GetPost` fonksiyonunu kullanarak blockchain state'inden belirtilen `Id`'ye sahip bir post nesnesi alır ve `found` boolean değişkeninin değerini kontrol ederek postun bulunup bulunmadığını kontrol eder. Gönderi bulunamazsa, `sdkerrors.ErrKeyNotFound` türünde bir hata döndürür.

Gönderi bulunmuşsa, fonksiyon alan olarak alınan gönderi nesnesine sahip yeni bir `types.QueryShowPostResponse` nesnesi oluşturur ve bu nesneye bir işaretçi ile `nil` hatası döndürür.

### `QueryShowPostResponse` öğesini değiştirme

Alanı bir işaretçi olmadan oluşturmak için `QueryShowPostResponse` mesajındaki gönderi alanına `[(gogoproto.nullable) = false]` seçeneğini ekleyin.

proto/blog/blog/query.proto

```
message QueryShowPostResponse {
  Post post = 1 [(gogoproto.nullable) = false];
}
```

Proto'dan Go dosyaları oluşturmak için komutu çalıştırın:

```
ignite generate proto-go
```
