# Blok Zincirler Arası İletişim: Temel Bilgiler

Blok Zincirler Arası İletişim protokolü (IBC) Cosmos SDK ekosisteminin önemli bir parçasıdır. Hello World eğitimi, bilgisayar programcılığında çok eski bir gelenektir. Bu eğitim, blok zinciri boyunca paketlerin nasıl oluşturulacağı ve gönderileceği konusunda bir anlayış oluşturur. Bu temel bilgi, Cosmos SDK ile blok zincirleri arasında gezinmenize yardımcı olur.

Şunları nasıl yapacağınızı öğreneceksiniz:

* Blok zincirleri arasında paketler oluşturmak ve göndermek için IBC'yi kullanın.
* Cosmos SDK ve Ignite CLI Relayer kullanarak blok zincirleri arasında gezinin.
* Temel bir blog gönderisi oluşturun ve gönderiyi başka bir blok zincirine kaydedin.

### IBC nedir?

Blok Zincirleri Arası İletişim protokolü (IBC) blok zincirlerinin birbirleriyle konuşmasını sağlar. IBC, farklı egemen blok zincirleri arasında aktarımı gerçekleştirir. Bu uçtan uca, bağlantı odaklı, durum bilgisi içeren protokol, heterojen blok zincirleri arasında güvenilir, sıralı ve kimliği doğrulanmış iletişim sağlar.

Cosmos SDK'daki IBC protokolü, iki blok zinciri arasındaki etkileşim için standarttır. IBCmodül arayüzü, paketlerin ve mesajların gönderen ve alan blok zinciri tarafından yorumlanmak üzere nasıl oluşturulduğunu tanımlar.

IBC aktarıcısı, IBC özellikli zincir kümeleri arasında bağlantı kurmanızı sağlar. Bu eğitim size iki blok zincirini nasıl oluşturacağınızı ve ardından iki blok zincirini bağlamak için Ignite CLI ile aktarıcıyı nasıl başlatacağınızı ve kullanacağınızı öğretir.

Bu eğitim modüller, IBC paketleri, aktarıcı ve IBC üzerinden yönlendirilen paketlerin yaşam döngüsü gibi temel konuları kapsamaktadır.

### Blockchain oluşturma

Hello World mesajını içeren diğer blok zincirlerine gönderiler yazmak için bir blog modülüne sahip bir blok zinciri uygulaması oluşturun. Bu eğitimde, Cosmos SDK evreni için Hello Mars, Hello Cosmos ve Hello Earth mesajlarını içeren gönderiler yazabilirsiniz.

Bu basit örnek için, başlık ve metin içeren bir gönderi işlemine sahip bir blog modülü içeren bir uygulama oluşturun.

Mantığı tanımladıktan sonra, bu modülün yüklü olduğu iki blok zinciri çalıştırın.

Zincirler IBC kullanarak birbirleri arasında gönderi gönderebilir.

Gönderen zincirde, onaylanan ve zaman aşımına uğrayan gönderileri kaydedin.

İşlem alıcı zincir tarafından onaylandıktan sonra, gönderinin her iki blok zincirine de kaydedildiğini bilirsiniz.

Gönderen zincir postID ek verisine sahiptir.

Onaylanan ve zaman aşımına uğrayan gönderiler, gönderinin başlığını ve hedef zincirini içerir. Bu tanımlayıcılar

parametre zincirinde görülebilir. Aşağıdaki grafik IBC'den geçen bir paketin yaşam döngüsünü göstermektedir.

<figure><img src="https://docs.ignite.com/assets/images/packet_sendpost-250db83c05d6472196790d0f04514173.png" alt=""><figcaption></figcaption></figure>

### Blockchain uygulamanızı oluşturun

Blockchain uygulamasını ve blog modülünü iskelelemek için Ignite CLI kullanın.

#### Yeni bir blok zinciri oluşturun

Planet adında yeni bir blok zincirinin iskeletini oluşturmak:

```
ignite scaffold chain planet --no-module
cd planet
```

Ev dizininizde planet adında yeni bir dizin oluşturulur. Planet dizini çalışan bir blockchain uygulaması içerir.

Blog modülünü blockchain'inizin içinde iskeleleyin

Ardından, Ignite CLI kullanarak IBC özelliklerine sahip bir blog modülünün iskelesini oluşturun. Blog modülü, blog gönderileri oluşturma ve bunları IBC aracılığıyla ikinci blockchain'e yönlendirme mantığını içerir.

`Blog` adlı bir modülü iskelelemek için:

```
ignite scaffold module blog --ibc
```

Bir IBC modülünün kodunu içeren yeni bir dizin planet/x/blog içinde oluşturulur. Ibc bayrağıyla iskelelenen modüller, iskelelenen IBC modülünün tüm mantığını içerir.

Türler için CRUD eylemleri oluşturun Ardından, blog modülü türleri için CRUD eylemlerini oluşturun.

Oluşturma, okuma, güncelleme ve silme (CRUD) eylemlerine yönelik şablon kodunu iskelelemek için ignite scaffold list komutunu kullanın.

Bu ignite iskele listesi komutları aşağıdaki işlemler için CRUD kodu oluşturur:

* Blog gönderileri oluşturma
* ```
  ignite scaffold list post title content creator --no-message --module blog
  ```
*   Gönderilen gönderiler için onayları işleme

    ```
    ignite scaffold list sentPost postID title chain creator --no-message --module blog
    ```
*   Gönderi zaman aşımlarını yönetme

    ```
    ignite scaffold list timedoutPost title chain creator --no-message --module blog
    ```

İskele kodu, veri yapılarını, mesajları, mesaj işleyicilerini, durumu değiştirmek için tutucuları ve CLI komutlarını tanımlamak için proto dosyalarını içerir.

#### Ignite CLI İskele Listesi Komutlarına Genel Bakış

```
ignite scaffold list [typeName] [field1] [field2] ... [flags]
```

ignite scaffold list \[typeName] komutunun ilk bağımsız değişkeni, oluşturulmakta olan türün adını belirtir. Blog uygulaması için post, sentPost ve timedoutPost türlerini oluşturdunuz.

Sonraki bağımsız değişkenler, türle ilişkilendirilen alanları tanımlar. Blog uygulaması için title, content, postID ve chain alanlarını oluşturdunuz.

\--module bayrağı, yeni işlem türünün hangi modüle ekleneceğini tanımlar. Bu isteğe bağlı bayrak, Ignite CLI uygulamanızda birden fazla modülü yönetmenizi sağlar. Bayrak mevcut olmadığında, tür, deponun adıyla eşleşen modülde iskele haline getirilir.

Yeni bir tür iskelelendiğinde, varsayılan davranış, CRUD işlemleri için kullanıcılar tarafından gönderilebilecek mesajları iskelelemektir. no-message bayrağı bu özelliği devre dışı bırakır. Gönderilerin IBC paketlerinin alınması üzerine oluşturulmasını ve doğrudan bir kullanıcının mesajlarından oluşturulmamasını istediğiniz için uygulama için mesajlar seçeneğini devre dışı bırakın.

Gönderilebilir ve yorumlanabilir bir IBC paketini iskeletleyin Blog gönderisinin başlığını ve içeriğini içeren bir paket için kod oluşturmalısınız.

Ignite packet komutu, başka bir blok zincirine gönderilebilecek bir IBC paketinin mantığını oluşturur.

Başlık ve içerik hedef zincirde saklanır.

Gönderen zincirde postID onaylanır.

Gönderilebilir ve yorumlanabilir bir IBC paketinin iskeletini oluşturmak için:

```
ignite scaffold packet ibcPost title content --ack postID --module blog
```

ibcPost paketindeki alanların daha önce oluşturduğunuz gönderi türündeki alanlarla eşleştiğine dikkat edin.

* ack bayrağı, gönderen blok zincirine hangi tanımlayıcının döndürüleceğini tanımlar.
* Modül bayrağı, paketin belirli bir IBC modülünde oluşturulacağını belirtir.

ignite packet komutu ayrıca bir IBC paketi gönderebilen CLI komutunun da iskelesini oluşturur:

```
planetd tx blog send-ibcPost [portID] [channelID] [title] [content]
```

### Kaynak kodunu değiştirin

Türleri ve işlemleri oluşturduktan sonra, veritabanındaki güncellemeleri yönetmek için mantığı manuel olarak eklemeniz gerekir. Bu eğitimde daha önce belirtildiği gibi verileri kaydetmek için kaynak kodunu değiştirin.

#### Blog yazısı paketine içerik oluşturucu ekleyin

IBC paketinin yapısını tanımlayan proto dosyası ile başlayın.

Alıcı blockchain'de gönderinin yaratıcısını tanımlamak için, paketin içine creator alanını ekleyin. Bu alan doğrudan komutta belirtilmemiştir çünkü SendIbcPost CLI komutunda otomatik olarak bir parametre haline gelecektir.

proto/planet/blog/packet.proto

```
message IbcPostPacketData {
  string title = 1;
  string content = 2;
  string creator = 3;
}
```

Alıcı zincirin bir blog gönderisinin oluşturucusu hakkında içeriğe sahip olduğundan emin olmak için, IBC paketine msg.Creator değerini ekleyin.

* Mesajı gönderenin içeriği otomatik olarak SendIbcPost mesajına dahil edilir.
* Gönderen, iletiyi imzalayan kişi olarak doğrulanır, bu nedenle msg.Sender değerini yeni pakete oluşturucu olarak ekleyebilirsiniz
* IBC üzerinden gönderilmeden önce.

x/blog/keeper/msg\_server\_ibc\_post.go

```
package keeper

import (
    // ...
    "planet/x/blog/types"
)

func (k msgServer) SendIbcPost(goCtx context.Context, msg *types.MsgSendIbcPost) (*types.MsgSendIbcPostResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: logic before transmitting the packet

    // Construct the packet
    var packet types.IbcPostPacketData

    packet.Title = msg.Title
    packet.Content = msg.Content
    packet.Creator = msg.Creator

    // Transmit the packet
    err := k.TransmitIbcPostPacket(
        ctx,
        packet,
        msg.Port,
        msg.ChannelID,
        clienttypes.ZeroHeight(),
        msg.TimeoutTimestamp,
    )
    if err != nil {
        return nil, err
    }

    return &types.MsgSendIbcPostResponse{}, nil
}
```

Gönderiyi alın

Birincil işlem mantığı için yöntemler `x/blog/keeper/ibc_post.go` dosyasındadır. IBC paketlerini yönetmek için bu yöntemleri kullanın:

* `TransmitIbcPostPacket`, paketi IBC üzerinden göndermek için manuel olarak çağrılır. Bu yöntem ayrıca paket IBC üzerinden başka bir blockchain uygulamasına gönderilmeden önceki mantığı da tanımlar.
* `OnRecvIbcPostPacket` kancası, zincir üzerinde bir paket alındığında otomatik olarak çağrılır. Bu yöntem paket alım mantığını tanımlar.
* `OnAcknowledgementIbcPostPacket` kancası, gönderilen bir paket kaynak zincirde onaylandığında çağrılır. Bu yöntem, paket alındığında mantığı tanımlar.
* `OnTimeoutIbcPostPacket` kancası, gönderilen bir paket zaman aşımına uğradığında çağrılır. Bu yöntem, paket hedef zincirde alınmadığında mantığı tanımlar

Veri tablolarının uygun şekilde değiştirilmesi için bu işlevlerin içine mantık eklemek üzere kaynak kodunu değiştirmeniz gerekir.

Gönderi mesajının alınması üzerine, alıcı zincirde başlık ve içerik ile yeni bir gönderi oluşturun.

Bir mesajın kaynaklandığı blok zinciri uygulamasını ve mesajı kimin oluşturduğunu tanımlamak için aşağıdaki formatta bir tanımlayıcı kullanın:

```
<portID>-<channelID>-<creatorAddress>
```

Son olarak, Ignite CLI tarafından oluşturulan AppendPost işlevi, eklenen yeni gönderinin kimliğini döndürür. Bu değeri onaylama yoluyla kaynak zincirine döndürebilirsiniz.

Paketi aldığınızda tür örneğini `PostID` olarak ekleyin:

* Context `ctx`, işlemden başlık verilerine sahip [değişmez bir veri yapısıdır](https://docs.cosmos.network/main/core/context.html#go-context-package). [Context'in nasıl başlatıldığını](https://github.com/cosmos/cosmos-sdk/blob/main/types/context.go#L71) görün.
* Daha önce tanımladığınız tanımlayıcı biçimi
* `title`, blog gönderisinin başlığıdır
* `content` blog yazısının içeriğidir

`x/blog/keeper/ibc_post.go` dosyasında, `"strconv"`u `"errors"`in altına aktardığınızdan emin olun:

x/blog/keeper/ibc\_post.go

```
import (
    //...

    "strconv"

// ...
)
```

Ardından `OnRecvIbcPostPacket` keeper işlevini aşağıdaki kodla değiştirin:

```
package keeper

// ...

func (k Keeper) OnRecvIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) (packetAck types.IbcPostPacketAck, err error) {
    // validate packet data upon receiving
    if err := data.ValidateBasic(); err != nil {
        return packetAck, err
    }

    id := k.AppendPost(
        ctx,
        types.Post{
            Creator: packet.SourcePort + "-" + packet.SourceChannel + "-" + data.Creator,
            Title:   data.Title,
            Content: data.Content,
        },
    )

    packetAck.PostID = strconv.FormatUint(id, 10)

    return packetAck, nil
}
```

#### Gönderi onayını alın

Gönderen blok zincirinde bir `sentPost` saklayın, böylece gönderinin hedef zincirde alındığını bilirsiniz.

Gönderiyi tanımlamak için başlığı ve hedefi saklayın.

Bir paket iskeletlendiğinde, alınan onay verileri için varsayılan tür, paket işleminin başarısız olup olmadığını tanımlayan bir türdür. `OnRecvIbcPostPacket` paketten bir hata döndürürse `Acknowledgement_Error` tipi ayarlanır.

x/blog/keeper/ibc\_post.go

```
package keeper

// ...

// x/blog/keeper/ibc_post.go
func (k Keeper) OnAcknowledgementIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData, ack channeltypes.Acknowledgement) error {
    switch dispatchedAck := ack.Response.(type) {
    case *channeltypes.Acknowledgement_Error:
        // We will not treat acknowledgment error in this tutorial
        return nil
    case *channeltypes.Acknowledgement_Result:
        // Decode the packet acknowledgment
        var packetAck types.IbcPostPacketAck

        if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
            // The counter-party module doesn't implement the correct acknowledgment format
            return errors.New("cannot unmarshal acknowledgment")
        }

        k.AppendSentPost(
            ctx,
            types.SentPost{
                Creator: data.Creator,
                PostID:  packetAck.PostID,
                Title:   data.Title,
                Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
            },
        )

        return nil
    default:
        return errors.New("the counter-party module does not implement the correct acknowledgment format")
    }
}
```

#### Zaman aşımına uğrayan paketle ilgili bilgileri saklar

Hedef zincirler tarafından alınmamış gönderileri `timedoutPost` gönderilerinde saklayın. Bu mantık `sentPost` ile aynı formatı izler.

x/blog/keeper/ibc\_post.go

```
func (k Keeper) OnTimeoutIbcPostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IbcPostPacketData) error {
    k.AppendTimedoutPost(
        ctx,
        types.TimedoutPost{
            Creator: data.Creator,
            Title:   data.Title,
            Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
        },
    )

    return nil
}

```

Bu son adım temel `blog` modülü kurulumunu tamamlar. Blok zinciri artık hazır!

### IBC Modüllerini test edin <a href="#use-the-ibc-modules" id="use-the-ibc-modules"></a>

Artık blok zincirini çalıştırabilir ve bir blok zinciri uygulamasından diğerine bir blog gönderisi gönderebilirsiniz. Bu sonraki adımları tamamlamak için birden fazla terminal penceresi gereklidir.

#### IBC modüllerini test edin

IBC modülünü test etmek için aynı makinede iki blok zinciri ağı başlatın. Her iki blok zinciri de aynı kaynak kodunu kullanır. Her blok zincirinin benzersiz bir zincir kimliği vardır.

Bir blockchain `earth` ve diğer blok zinciri `mars` olarak adlandırılır.

Proje dizininde `earth.yml` ve `mars.yml` dosyaları gereklidir:

earth.yml

```
version: 1
build:
  proto:
    path: proto
    third_party_paths:
    - third_party/proto
    - proto_vendor
accounts:
- name: alice
  coins:
  - 1000token
  - 100000000stake
- name: bob
  coins:
  - 500token
  - 100000000stake
faucet:
  name: bob
  coins:
  - 5token
  - 100000stake
  host: 0.0.0.0:4500
genesis:
  chain_id: earth
validators:
- name: alice
  bonded: 100000000stake
  home: $HOME/.earth
```

mars.yml

```
version: 1
build:
  proto:
    path: proto
    third_party_paths:
    - third_party/proto
    - proto_vendor
accounts:
- name: alice
  coins:
  - 1000token
  - 1000000000stake
- name: bob
  coins:
  - 500token
  - 100000000stake
faucet:
  name: bob
  coins:
  - 5token
  - 100000stake
  host: :4501
genesis:
  chain_id: mars
validators:
- name: alice
  bonded: 100000000stake
  app:
    api:
      address: :1318
    grpc:
      address: :9092
    grpc-web:
      address: :9093
  config:
    p2p:
      laddr: :26658
    rpc:
      laddr: :26659
      pprof_laddr: :6061
  home: $HOME/.mars
```

Bir terminal penceresi açın ve `earth` blockchain'i başlatmak için aşağıdaki komutu çalıştırın:

```
ignite chain serve -c earth.yml
```

Farklı bir terminal penceresi açın ve `mars` blockchain'ini başlatmak için aşağıdaki komutu çalıştırın:

```
ignite chain serve -c mars.yml
```

#### Mevcut Relayer ve Ignite CLI Yapılandırmalarını Kaldırma

Daha önce relayer kullandıysanız, çıkan relayer ve Ignite CLI konfigürasyonlarını kaldırmak için aşağıdaki adımları izleyin:

* Blok zincirlerinizi durdurun ve önceki yapılandırma dosyalarını silin:

```
rm -rf ~/.ignite/relayer
```

Mevcut aktarıcı yapılandırmaları yoksa, komut hiçbir eşleşme bulunamadı sonucunu döndürür ve hiçbir işlem yapılmaz.

#### Aktarıcıyı yapılandırma ve başlatma

İlk olarak, aktarıcıyı yapılandırın. Ignite CLI configure komutunu `--advanced` seçeneği ile kullanın:

```
ignite relayer configure -a \
  --source-rpc "http://0.0.0.0:26657" \
  --source-faucet "http://0.0.0.0:4500" \
  --source-port "blog" \
  --source-version "blog-1" \
  --source-gasprice "0.0000025stake" \
  --source-prefix "cosmos" \
  --source-gaslimit 300000 \
  --target-rpc "http://0.0.0.0:26659" \
  --target-faucet "http://0.0.0.0:4501" \
  --target-port "blog" \
  --target-version "blog-1" \
  --target-gasprice "0.0000025stake" \
  --target-prefix "cosmos" \
  --target-gaslimit 300000
```

İstendiğinde, `Source Account` ve `Target Account` için varsayılan değerleri kabul etmek üzere Enter tuşuna basın.

Çıktı aşağıdaki gibi görünür:

```
---------------------------------------------
Setting up chains
---------------------------------------------

🔐  Account on "source" is "cosmos1xcxgzq75yrxzd0tu2kwmwajv7j550dkj7m00za"

 |· received coins from a faucet
 |· (balance: 100000stake,5token)

🔐  Account on "target" is "cosmos1nxg8e4mfp5v7sea6ez23a65rvy0j59kayqr8cx"

 |· received coins from a faucet
 |· (balance: 100000stake,5token)

⛓  Configured chains: earth-mars
```

Yeni bir terminal penceresinde relayer işlemini başlatın:

```
ignite relayer connect
```

Sonuçlar:

```
------
Paths
------

earth-mars:
    earth > (port: blog) (channel: channel-0)
    mars  > (port: blog) (channel: channel-0)

------
Listening and relaying packets between chains...
------
```

#### Paketleri gönder

Artık paket gönderebilir ve alınan gönderileri doğrulayabilirsiniz:

```
planetd tx blog send-ibc-post blog channel-0 "Hello" "Hello Mars, I'm Alice from Earth" --from alice --chain-id earth --home ~/.earth
```

Gönderinin Mars'a ulaştığını doğrulamak için:

```
planetd q blog list-post --node tcp://localhost:26659
```

Paket alındı:

```
Post:
  - content: Hello Mars, I'm Alice from Earth
    creator: blog-channel-0-cosmos1aew8dk9cs3uzzgeldatgzvm5ca2k4m98xhy20x
    id: "0"
    title: Hello
pagination:
  next_key: null
  total: "1"
```

Paketin Dünya'da onaylanıp onaylanmadığını kontrol etmek için:

```
planetd q blog list-sent-post
```

Çıktı:

```
SentPost:
  - chain: blog-channel-0
    creator: cosmos1aew8dk9cs3uzzgeldatgzvm5ca2k4m98xhy20x
    id: "0"
    postID: "0"
    title: Hello
pagination:
  next_key: null
  total: "1"
```

Zaman aşımını test etmek için, bir paketin zaman aşımı süresini 1 nanosaniyeye ayarlayın, paketin zaman aşımına uğradığını doğrulayın ve zaman aşımına uğrayan gönderileri kontrol edin:

```
planetd tx blog send-ibc-post blog channel-0 "Sorry" "Sorry Mars, you will never see this post" --from alice --chain-id earth --home ~/.earth --packet-timeout-timestamp 1
```

Zaman aşımına uğramış gönderileri kontrol edin:

```
planetd q blog list-timedout-post
```

Sonuçlar:

```
TimedoutPost:
  - chain: blog-channel-0
    creator: cosmos1fhpcsxn0g8uask73xpcgwxlfxtuunn3ey5ptjv
    id: "0"
    title: Sorry
pagination:
  next_key: null
  total: "2"
```

Mars'tan da posta gönderebilirsiniz:

```
planetd tx blog send-ibc-post blog channel-0 "Hello" "Hello Earth, I'm Alice from Mars" --from alice --chain-id mars --home ~/.mars --node tcp://localhost:26659
```

Dünya'daki liste gönderisi:

```
planetd q blog list-post
```

Sonuçlar:

```
Post:
  - content: Hello Earth, I'm Alice from Mars
    creator: blog-channel-0-cosmos1xtpx43l826348s59au24p22pxg6q248638q2tf
    id: "0"
    title: Hello
pagination:
  next_key: null
  total: "1"
```

### Tebrikler 🎉

Bu eğitimi tamamlayarak Blok Zincirler Arası İletişim protokolünü (IBC) kullanmayı öğrendiniz.

İşte bu eğitimde başardıklarınız:

* IBC modülleri olarak iki Hello blockchain uygulaması oluşturdunuz
* CRUD eylem mantığını eklemek için oluşturulan kod değiştirildi
* İki blok zincirini birbirine bağlamak için Ignite CLI aktarıcısını yapılandırdı ve kullandı
* IBC paketlerinin bir blok zincirinden diğerine aktarılması
