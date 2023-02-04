# Gönderileri silme

Bu bölümde, bir "gönderiyi sil" mesajını işleme sürecine odaklanacağız.

### Gönderileri kaldırma

x/blog/keeper/post.go

```
func (k Keeper) RemovePost(ctx sdk.Context, id uint64) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
    store.Delete(GetPostIDBytes(id))
}
```

`RemovePost` fonksiyonu iki argüman alır: bir context nesnesi `ctx` ve bir işaretsiz tamsayı `id`. Fonksiyon, verilen `id` ile ilişkili anahtar-değer çiftini silerek bir anahtar-değer deposundan bir gönderiyi kaldırır. Anahtar-değer deposuna, bağlamın anahtar-değer deposunu ve `PostKey` sabitine dayalı bir öneki kullanarak yeni bir store oluşturmak için `prefix` paketi kullanılarak oluşturulan `store` değişkeni kullanılarak erişilir. Daha sonra Delete yöntemi, `id`'yi silinecek anahtar olarak bir byte dilimine dönüştürmek için `GetPostIDBytes` işlevi kullanılarak store nesnesi üzerinde çağrılır.

### Gönderileri silme

x/blog/keeper/msg\_server\_delete\_post.go

```
package keeper

import (
    "context"
    "fmt"

    "blog/x/blog/types"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) DeletePost(goCtx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    val, found := k.GetPost(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
    }
    if msg.Creator != val.Creator {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }
    k.RemovePost(ctx, msg.Id)
    return &types.MsgDeletePostResponse{}, nil
}
```

DeletePost iki bağımsız değişken alır: context.Context türünde bir context goCtx ve \*types.MsgDeletePost türünde bir mesaja işaretçi. Fonksiyon \*types.MsgDeletePostResponse türünde bir mesaja bir işaretçi ve bir hata döndürür.

Fonksiyonun içinde, `sdk.UnwrapSDKContext` fonksiyonu kullanılarak context açılır ve `GetPost` fonksiyonu kullanılarak mesajda belirtilen ID'ye sahip gönderinin değeri alınır. Gönderi bulunamazsa, `sdkerrors.Wrap` işlevi kullanılarak bir hata döndürülür. İletinin oluşturucusu gönderinin oluşturucusuyla eşleşmiyorsa, başka bir hata döndürülür. Bu kontrollerin her ikisi de başarılı olursa, gönderiyi silmek için `RemovePost` fonksiyonu bağlam ve gönderinin ID'si ile çağrılır. Son olarak, fonksiyon veri içermeyen bir yanıt mesajı ve `nil` hatası döndürür.

Kısacası, `DeletePost` bir gönderiyi silme talebini ele alır ve silmeden önce talep edenin gönderinin yaratıcısı olduğundan emin olur.

### Özet

Keeper paketindeki `RemovePost` ve `DeletePost` yöntemlerinin uygulamasını tamamladığınız için tebrikler! Bu yöntemler, sırasıyla bir gönderiyi depodan kaldırmak ve bir gönderiyi silme isteğini ele almak için işlevsellik sağlar.
