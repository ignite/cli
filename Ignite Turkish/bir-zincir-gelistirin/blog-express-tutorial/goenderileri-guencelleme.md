# Gönderileri güncelleme

Bu bölümde, bir "gönderiyi güncelle" mesajını işleme sürecine odaklanacağız.

Bir gönderiyi güncellemek için, "Get" işlemini kullanarak belirli bir gönderiyi mağazadan almanız, değerleri değiştirmeniz ve ardından "Set" işlemini kullanarak güncellenmiş gönderiyi mağazaya geri yazmanız gerekir.

Önce bir getter ve setter mantığı uygulayalım.

### Gönderileri alma

`post.go` içinde `GetPost` keeper yöntemini uygulayın:

x/blog/keeper/post.go

```
func (k Keeper) GetPost(ctx sdk.Context, id uint64) (val types.Post, found bool) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
    b := store.Get(GetPostIDBytes(id))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

`GetPost` iki argüman alır: bir context `ctx` ve alınacak gönderinin kimliğini temsil eden `uint64` türünde bir `id`. Gönderinin değerlerini içeren `types.Post` yapısını ve gönderinin veritabanında bulunup bulunmadığını gösteren bir boolean değerini döndürür.

Fonksiyon ilk olarak `prefix.NewStore` yöntemini kullanarak, context'ten anahtar-değer deposunu ve `types.PostKey` sabitine uygulanan `types.KeyPrefix` fonksiyonunu argüman olarak geçirerek bir mağaza oluşturur. Daha sonra `store.Get` yöntemini kullanarak mağazadan gönderiyi almaya çalışır ve gönderinin kimliğini bir byte dilimi olarak iletir. Gönderi mağazada bulunamazsa, boş bir `types.Post` yapısı ve false boolean değeri döndürür.

Gönderi depoda bulunursa, fonksiyon `cdc.MustUnmarshal` yöntemini kullanarak alınan byte dilimini bir `types.Post` yapısına ayırır ve argüman olarak val değişkenine bir işaretçi iletir. Ardından val yapısını ve gönderinin veritabanında bulunduğunu gösteren true boolean değerini döndürür.

### Gönderileri ayarlama

`SetPost` keeper yöntemini `post.go`'da uygulayın:

x/blog/keeper/post.go

```
func (k Keeper) SetPost(ctx sdk.Context, post types.Post) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
    b := k.cdc.MustMarshal(&post)
    store.Set(GetPostIDBytes(post.Id), b)
}
```

`SetPost` iki argüman alır: bir context `ctx` ve post için güncellenmiş değerleri içeren bir `types.Post` struct. Fonksiyon hiçbir şey döndürmez.

Fonksiyon ilk olarak `prefix.NewStore` yöntemini kullanarak bir depo oluşturur, context'ten anahtar-değer deposunu ve `types.PostKey` sabitine uygulanan `types.KeyPrefix` fonksiyonunu argüman olarak iletir. Daha sonra `cdc.MustMarshal` yöntemini kullanarak güncellenmiş post yapısını bir byte dilimi halinde marshall eder ve argüman olarak post yapısına bir işaretçi iletir. Son olarak, `store.Set` yöntemini kullanarak depodaki gönderiyi günceller, gönderinin kimliğini bir byte dilimi olarak ve harmanlanmış gönderi yapısını argüman olarak iletir.

### Gönderileri güncelleme

x/blog/keeper/msg\_server\_update\_post.go

```
package keeper

import (
    "context"
    "fmt"

    "blog/x/blog/types"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdatePost(goCtx context.Context, msg *types.MsgUpdatePost) (*types.MsgUpdatePostResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    var post = types.Post{
        Creator: msg.Creator,
        Id:      msg.Id,
        Title:   msg.Title,
        Body:    msg.Body,
    }
    val, found := k.GetPost(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
    }
    if msg.Creator != val.Creator {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }
    k.SetPost(ctx, post)
    return &types.MsgUpdatePostResponse{}, nil
}
```

`UpdatePost`, girdi olarak bir context ve `MsgUpdatePost` mesajını alır ve `MsgUpdatePostResponse` yanıtı ile bir hata döndürür. Fonksiyon ilk olarak sağlanan `msg.Id`'yi kullanarak veritabanından gönderinin geçerli değerlerini alır ve gönderinin var olup olmadığını ve `msg.Creator`'ın gönderinin geçerli sahibiyle aynı olup olmadığını kontrol eder. Bu kontrollerden herhangi biri başarısız olursa hata döndürür. Her iki kontrol de başarılı olursa, veritabanındaki gönderiyi msg'de sağlanan yeni değerlerle günceller ve hata olmadan bir yanıt döndürür.

### Özet

Tebrikler! Bir mağaza içindeki gönderileri yönetmek için bir dizi önemli yöntemi başarıyla uyguladınız.

`GetPost` yöntemi, benzersiz kimlik numarasına veya gönderi kimliğine göre mağazadan belirli bir gönderiyi almanıza olanak tanır. Bu, belirli bir gönderiyi bir kullanıcıya görüntülemek veya güncellemek için yararlı olabilir.

SetPost yöntemi, mağazadaki mevcut bir gönderiyi güncellemenizi sağlar. Bu, hataları düzeltmek veya yeni bilgiler elde edildikçe bir gönderinin içeriğini güncellemek için yararlı olabilir.

Son olarak, blok zinciri bir gönderide güncelleme talep eden bir mesajı her işlediğinde çağrılan `UpdatePost` yöntemini uyguladınız.
