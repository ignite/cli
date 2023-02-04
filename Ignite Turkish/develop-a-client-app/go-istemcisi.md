# Go istemcisi

### Go programlama dilinde bir istemci

Bu eğitimde, bir blockchain için istemci olarak hizmet veren bağımsız bir Go programının nasıl oluşturulacağını göstereceğiz. Standart bir blockchain kurmak için Ignite CLI kullanacağız. Blockchain ile iletişim kurmak için, blockchain ile etkileşim için kullanımı kolay bir arayüz sağlayan `cosmosclient`paketini kullanacağız. İşlem göndermek ve blockchaini sorgulamak için `cosmosclient`paketini nasıl kullanacağınızı öğreneceksiniz. Bu eğitimin sonunda, Go ve `cosmosclient` paketini kullanarak bir blockchain için bir istemcinin nasıl oluşturulacağını iyi bir şekilde anlayacaksınız.

Ignite CLI kullanarak bir blockchain oluşturmak için aşağıdaki komutu kullanın:

```
ignite scaffold chain blog
```

Bu, "blog" adında yeni bir Cosmos SDK blockchain'i oluşturacaktır.

Blockchain oluşturulduktan sonra, blog gönderileri üzerinde oluşturma, okuma, güncelleme ve silme (CRUD) işlemlerini gerçekleştirmenizi sağlayacak bir "blog" modeli için kod oluşturabilirsiniz. Bunu yapmak için aşağıdaki komutu kullanabilirsiniz:

```
cd blog
ignite scaffold list post title body
```

Bu, blog gönderilerini oluşturma, okuma, güncelleme ve silme işlevleri de dahil olmak üzere "blog" modeli için gerekli kodu oluşturacaktır. Bu kod sayesinde artık blog gönderileri üzerinde CRUD işlemleri gerçekleştirmek için blockchaininizi kullanabilirsiniz. Oluşturulan kodu yeni blog gönderileri oluşturmak, mevcut olanları almak, içeriklerini güncellemek ve gerektiğinde silmek için kullanabilirsiniz. Bu size blog gönderilerini yönetme becerisine sahip tamamen işlevsel bir Cosmos SDK blok zinciri sağlayacaktır.

Blockchain node'unuzu aşağıdaki komutla başlatın:

`blog` dizini ile aynı seviyede `blogclient` adında yeni bir dizin oluşturun. Adından da anlaşılacağı gibi `blogclient`, `blog` blockchain'iniz için bir istemci görevi gören bağımsız bir Go programı içerecektir.

Bu komut bulunduğunuz konumda `blogclient` adında yeni bir dizin oluşturacaktır. Terminal pencerenize `ls` yazarsanız, hem `blog` hem de `blogclient` dizinlerinin listelendiğini görürsünüz.

Blogclient dizini içinde yeni bir Go paketi başlatmak için aşağıdaki komutu kullanabilirsiniz:

```
cd blogclient
go mod init blogclient
```

Bu, `blogclient` dizininde paket ve kullanılan Go sürümü hakkında bilgi içeren bir `go.mod` dosyası oluşturacaktır.

Paketinizin bağımlılıklarını içe aktarmak için `go.mod` dosyasına aşağıdaki kodu ekleyebilirsiniz:

blogclient/go.mod

```
module blogclient

go 1.19

require (
    blog v0.0.0-00010101000000-000000000000
    github.com/ignite/cli v0.25.2
)

replace blog => ../blog
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
```

Paketiniz iki bağımlılığı içe aktaracaktır:

* mesaj types (tiplerini) ve bir sorgu istemcisini içeren `blog`
* `cosmosclient` paketi için `ignite`

`replace`yönergesi yerel `blog`dizinindeki paketi kullanır ve `blogclient`dizinine göreli bir yol olarak belirtilir.

Cosmos SDK `protobuf`paketinin özel bir sürümünü kullanır, bu nedenle doğru bağımlılığı belirtmek için `replace`yönergesini kullanın.

Son olarak, `blog`istemciniz için bağımlılıkları yükleyin:

#### `Main.go`'da istemcinin ana mantığı[​](broken-reference) <a href="#main-logic-of-the-client-in-maingo" id="main-logic-of-the-client-in-maingo"></a>

`blogclient` dizini içinde bir `main.go` dosyası oluşturun ve aşağıdaki kodu ekleyin:

blogclient/main.go

```
package main

import (
    "context"
    "fmt"
    "log"

    // Importing the general purpose Cosmos blockchain client
    "github.com/ignite/cli/ignite/pkg/cosmosclient"

    // Importing the types package of your blog blockchain
    "blog/x/blog/types"
)

func main() {
    ctx := context.Background()
    addressPrefix := "cosmos"

    // Create a Cosmos client instance
    client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix))
    if err != nil {
        log.Fatal(err)
    }

    // Account `alice` was initialized during `ignite chain serve`
    accountName := "alice"

    // Get account from the keyring
    account, err := client.Account(accountName)
    if err != nil {
        log.Fatal(err)
    }

    addr, err := account.Address(addressPrefix)
    if err != nil {
        log.Fatal(err)
    }

    // Define a message to create a post
    msg := &types.MsgCreatePost{
        Creator: addr,
        Title:   "Hello!",
        Body:    "This is the first post",
    }

    // Broadcast a transaction from account `alice` with the message
    // to create a post store response in txResp
    txResp, err := client.BroadcastTx(ctx, account, msg)
    if err != nil {
        log.Fatal(err)
    }

    // Print response from broadcasting a transaction
    fmt.Print("MsgCreatePost:\n\n")
    fmt.Println(txResp)

    // Instantiate a query client for your `blog` blockchain
    queryClient := types.NewQueryClient(client.Context())

    // Query the blockchain using the client's `PostAll` method
    // to get all posts store all posts in queryResp
    queryResp, err := queryClient.PostAll(ctx, &types.QueryAllPostRequest{})
    if err != nil {
        log.Fatal(err)
    }

    // Print response from querying all the posts
    fmt.Print("\n\nAll posts:\n\n")
    fmt.Println(queryResp)
}
```

Yukarıdaki kod, `blog` blockchain için bir istemci görevi gören bağımsız bir Go programı oluşturur. Genel amaçlı Cosmos blockchain istemcisi ve `blog` blockchain'inin `types` paketi dahil olmak üzere gerekli paketleri içe aktararak başlar.

`main` işlevde, kod bir Cosmos istemci örneği oluşturur ve adres önekini "cosmos" olarak ayarlar. Daha sonra anahtarlıktan `"alice"` adlı bir hesap alır ve adres önekini kullanarak hesabın adresini alır.

Daha sonra kod, "Merhaba!" başlıklı ve "Bu ilk yazı" gövdeli bir blog yazısı oluşturmak için bir mesaj tanımlar. Ardından, "alice" hesabından gönderiyi oluşturmak için mesajla birlikte bir işlem yayınlar ve yanıtı `txResp` değişkeninde saklar.

Kod daha sonra blog blockchain'i için bir sorgu istemcisi oluşturur ve tüm gönderileri almak üzere blockchain'i sorgulamak için bunu kullanır. Yanıtı `queryResp` değişkeninde saklar ve konsola yazdırır.

Son olarak kod, işlemin yayınlanmasından elde edilen yanıtı konsola yazdırır. Bu, kullanıcının istemciyi kullanarak `blog` blockchain'inde bir blog gönderisi oluşturma ve sorgulama sonuçlarını görmesini sağlar.

`cosmosclient` paketi hakkında daha fazla bilgi edinmek için [`cosmosclient` ](https://pkg.go.dev/github.com/ignite/cli/ignite/pkg/cosmosclient)için Go paket belgelerine başvurabilirsiniz. Bu belge, `Client` türünün `Options` ve `KeyringBackend` ile nasıl kullanılacağı hakkında bilgi sağlar.

Blog blok zincirinizin hala `ignite chain serve` ile çalıştığından emin olun.

Blockchain istemcisini çalıştırın:

Komut başarılı olursa, komutun çalıştırılmasının sonuçları terminale yazdırılacaktır. Çıktı, göz ardı edilebilecek bazı uyarılar içerebilir.

```
MsgCreatePost:

code: 0
codespace: ""
data: 12220A202F626C6F672E626C6F672E4D7367437265617465506F7374526573706F6E7365
events:
- attributes:
  - index: true
    key: ZmVl
    value: null
  - index: true
    key: ZmVlX3BheWVy
    value: Y29zbW9zMWR6ZW13NzZ3enQ3cDBnajd3MzQyN2E0eHg3MjRkejAzd3hnOGhk
  type: tx
- attributes:
  - index: true
    key: YWNjX3NlcQ==
    value: Y29zbW9zMWR6ZW13NzZ3enQ3cDBnajd3MzQyN2E0eHg3MjRkejAzd3hnOGhkLzE=
  type: tx
- attributes:
  - index: true
    key: c2lnbmF0dXJl
    value: UWZncUJCUFQvaWxWVzJwNUJNTngzcDlvRzVpSXp0elhXdE9yMHcwVE00OEtlSkRqR0FEdU9VNjJiY1ZRNVkxTHdEbXNuYUlsTmc3VE9uMnJ2ZWRHSlE9PQ==
  type: tx
- attributes:
  - index: true
    key: YWN0aW9u
    value: L2Jsb2cuYmxvZy5Nc2dDcmVhdGVQb3N0
  type: message
gas_used: "52085"
gas_wanted: "300000"
height: "20"
info: ""
logs:
- events:
  - attributes:
    - key: action
      value: /blog.blog.MsgCreatePost
    type: message
  log: ""
  msg_index: 0
raw_log: '[{"msg_index":0,"events":[{"type":"message","attributes":[{"key":"action","value":"/blog.blog.MsgCreatePost"}]}]}]'
timestamp: ""
tx: null
txhash: 4F53B75C18254F96EF159821DDD665E965DBB576A5AC2B94CE863EB62E33156A

All posts:

Post:<title:"Hello!" body:"This is the first post" creator:"cosmos1dzemw76wzt7p0gj7w3427a4xx724dz03wxg8hd" > pagination:<total:1 >
```

Gördüğünüz gibi istemci bir işlemi başarıyla yayınladı ve blog gönderileri için zinciri sorguladı.

Lütfen terminalinizdeki çıktıdaki bazı değerlerin (işlem hash'i ve blok yüksekliği gibi) yukarıdaki çıktıdan farklı olabileceğini unutmayın.

Yeni gönderiyi `blogd q blog list-post` komutunu kullanarak onaylayabilirsiniz:

```
Post:
- body: This is the first post
  creator: cosmos1dzemw76wzt7p0gj7w3427a4xx724dz03wxg8hd
  id: "0"
  title: Hello!
pagination:
  next_key: null
  total: "0"
```

Harika bir iş başardınız! Cosmos SDK blok zinciriniz için bir Go istemcisi oluşturma, bir işlem gönderme ve zinciri sorgulama sürecini başarıyla tamamladınız.
