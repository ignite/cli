# Hello, World! In-depth Öğretici

Bu eğitimde "Hello, World!" işlevselliğini sıfırdan uygulayacaksınız. Oluşturacağınız uygulamanın işlevselliği "Ekspres öğretici" bölümünde oluşturduğunuzla aynı olacaktır, ancak burada süreci daha derinlemesine anlamak için manuel olarak yapacaksınız.

Başlamak için, yeni bir merhaba blok zinciri ile başlayalım. Önceki bölümde yaptığınız değişiklikleri geri alabilir ya da Ignite kullanarak yeni bir blok zinciri oluşturabilirsiniz. Her iki durumda da, üzerinde çalışmanız için hazır olan boş bir blok zincirine sahip olacaksınız.

```
ignite scaffold chain hello
```

### `SayHello` RPC <a href="#sayhello-rpc" id="sayhello-rpc"></a>

Cosmos SDK blok zincirlerinde sorgular, protokol buffer dosyalarındaki bir `Query` hizmetinde uzak prosedür çağrıları (RPC'ler) olarak tanımlanır. Yeni bir sorgu eklemek için aşağıdaki kodu modülünüzün query.proto dosyasına ekleyebilirsiniz:

proto/hello/hello/query.proto

```
service Query {
    rpc SayHello(QuerySayHelloRequest) returns (QuerySayHelloResponse) {
        option (google.api.http).get = "/hello/hello/say_hello/{name}";
    }
}
```

RPC, `QuerySayHelloRequest` türünde bir istek argümanı kabul eder ve `QuerySayHelloResponse` türünde bir değer döndürür. Bu türleri tanımlamak için `query.proto` dosyasına aşağıdaki kodu ekleyebilirsiniz:

proto/hello/hello/query.proto

```
message QuerySayHelloRequest {
  string name = 1;
}

message QuerySayHelloResponse {
  string name = 1;
}
```

query.proto'da tanımlanan türleri kullanmak için protokol buffer dosyalarını Go kaynak koduna dönüştürmeniz gerekir. Bu, blok zincirini oluşturup başlatacak ve protokol buffer dosyalarından Go kaynak kodunu otomatik olarak üretecek olan ignite chain serve çalıştırılarak yapılabilir. Alternatif olarak, blok zincirini oluşturmadan ve başlatmadan yalnızca protokol buffer dosyalarından Go kaynak kodunu oluşturmak için ignite generate proto-go'yu çalıştırabilirsiniz.

### `SayHello` keeper yöntemi

Sorgu, istek ve yanıt türlerini query.proto dosyasında tanımladıktan sonra, kodunuzda sorgu için mantığı uygulamanız gerekecektir. Bu genellikle isteği işleyen ve uygun yanıtı döndüren bir fonksiyon yazmayı içerir. Aşağıdaki içeriğe sahip yeni bir `query_say_hello.go` dosyası oluşturun:

x/hello/keeper/query\_say\_hello.go

```
package keeper

import (
    "context"
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "hello/x/hello/types"
)

func (k Keeper) SayHello(goCtx context.Context, req *types.QuerySayHelloRequest) (*types.QuerySayHelloResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)
    // TODO: Process the query
    _ = ctx
    return &types.QuerySayHelloResponse{Name: fmt.Sprintf("hello %s", req.Name)}, nil
}
```

Bu kod, `QuerySayHelloRequest` türünde bir isteği kabul eden ve `QuerySayHelloResponse` türünde bir değer döndüren bir `SayHello` işlevi tanımlar. İşlev önce isteğin geçerli olup olmadığını kontrol eder ve ardından `%s` yer tutucusunun değeri olarak sağlanan adla birlikte yanıt iletisini döndürerek sorguyu işler. Sorguyu işlemek ve uygun yanıtı döndürmek için gerektiğinde işleve blok zincirinden veri alma veya karmaşık işlemler gerçekleştirme gibi ek mantık ekleyebilirsiniz.

### `CmdSayHello` komutu <a href="#cmdsayhello-command" id="cmdsayhello-command"></a>

Sorgu mantığını uyguladıktan sonra, sorguyu çağırabilmeleri ve yanıtı alabilmeleri için istemcilerin kullanımına sunmanız gerekecektir. Bu genellikle sorgunun blok zincirinin uygulama programlama arayüzüne (API) eklenmesini ve kullanıcıların sorguyu kolayca göndermesine ve yanıtı almasına olanak tanıyan bir komut satırı arayüzü (CLI) komutu sağlamayı içerir.

Sorgu için bir CLI komutu sağlamak için `query_say_hello.go` dosyasını oluşturabilir ve `SayHello` işlevini çağıran ve yanıtı konsola yazdıran bir `CmdSayHello` komutu uygulayabilirsiniz.

x/hello/client/cli/query\_say\_hello.go

```
package cli

import (
    "strconv"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/spf13/cobra"

    "hello/x/hello/types"
)

var _ = strconv.Itoa(0)

func CmdSayHello() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "say-hello [name]",
        Short: "Query say-hello",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) (err error) {
            reqName := args[0]
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }
            queryClient := types.NewQueryClient(clientCtx)
            params := &types.QuerySayHelloRequest{
                Name: reqName,
            }
            res, err := queryClient.SayHello(cmd.Context(), params)
            if err != nil {
                return err
            }
            return clientCtx.PrintProto(res)
        },
    }
    flags.AddQueryFlagsToCmd(cmd)
    return cmd
}
```

Kod bir `CmdSayHello` komutu tanımlamaktadır. Komut, Go'da komut satırı uygulamaları oluşturmak için popüler bir çerçeve olan `cobra` kütüphanesi kullanılarak tanımlanmıştır. Komut, argüman olarak bir isim kabul eder ve bunu `types.QueryClient`'tan `SayHello` fonksiyonuna aktarılan bir `QuerySayHelloRequest` struct'ı oluşturmak için kullanır. `SayHello` fonksiyonu, `say-hello` sorgusunu blok zincirine göndermek için kullanılır ve yanıt res değişkeninde saklanır.

`QuerySayHelloRequest` yapısı, sorgu için istek ve yanıt türlerini tanımlayan bir Protokol Buffer dosyası olan `query.proto` dosyasında tanımlanır. `QuerySayHelloRequest` yapısı, yanıt mesajına dahil edilecek adı sağlamak için kullanılan dize türünde bir Name alanı içerir.

Sorgu gönderildikten ve yanıt alındıktan sonra, kod yanıtı konsola yazdırmak için `clientCtx.PrintProto` işlevini kullanır. `clientCtx` değişkeni, istemcinin yapılandırması ve bağlantı bilgileri dahil olmak üzere istemci bağlamına erişim sağlayan `client.GetClientQueryContext` işlevi kullanılarak elde edilir. `PrintProto` fonksiyonu, verilerin verimli bir şekilde serileştirilmesine ve seriden çıkarılmasına olanak tanıyan Protokol Arabelleği formatını kullanarak yanıtı yazdırmak için kullanılır.

`flags.AddQueryFlagsToCmd` işlevi, komuta sorguyla ilgili bayraklar eklemek için kullanılır. Bu, kullanıcıların komutu çağırırken düğüm URL'si ve diğer sorgu parametreleri gibi ek seçenekler belirtmesine olanak tanır. Bu bayraklar sorguyu yapılandırmak için kullanılır ve `SayHello` fonksiyonuna gerekli bilgileri sağlayarak blok zincirine bağlanmasına ve sorguyu göndermesine olanak tanır.

`CmdSayHello` komutunu kullanıcıların kullanımına sunmak için bu komutu zincirin ikili dosyasına eklemeniz gerekir. Bu genellikle `x/hello/client/cli/query.go` dosyasını değiştirerek ve `cmd.AddCommand(CmdSayHello())` ifadesini ekleyerek yapılır. Bu, `CmdSayHello` komutunu kullanılabilir komutlar listesine ekleyerek kullanıcıların komut satırı arayüzünden (CLI) çağırmasına olanak tanır.

x/hello/client/cli/query.go

```
func GetQueryCmd(queryRoute string) *cobra.Command {
    cmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }
    cmd.AddCommand(CmdQueryParams())
    cmd.AddCommand(CmdSayHello())
    return cmd
}
```

Bir CLI komutu sağladıktan sonra, kullanıcılar `say-hello` sorgusunu çağırabilecek ve uygun yanıtı alabileceklerdir.

Projenizin kaynak kodunda yaptığınız tüm değişiklikleri kaydedin ve bir blockchain node'u başlatmak için aşağıdaki komutu çalıştırın:

```
ignite chain serve
```

Sorguyu göndermek ve yanıtı almak için aşağıdaki komutu kullanın:

```
hellod q hello say-hello bob
```

Bu komut, blok zincirine "bob" adıyla bir "say-hello" sorgusu gönderecek ve "Hello, bob!" yanıtını konsola yazdıracaktır. Özel gereksinimlerinize uyması ve istediğiniz işlevselliği sağlaması için sorguyu ve yanıtı gerektiği gibi değiştirebilirsiniz.

"Hello, World!" eğitimini tamamladığınız için tebrikler! Bu eğitimde, bir protokol buffer dosyasında yeni bir sorguyu nasıl tanımlayacağınızı, kodunuzda sorgu mantığını nasıl uygulayacağınızı ve sorguyu blockchain'in API'si ve CLI aracılığıyla istemcilerin kullanımına nasıl sunacağınızı öğrendiniz. Eğitimde özetlenen adımları izleyerek, blok zincirinizden veri almak veya gerektiğinde diğer işlemleri gerçekleştirmek için kullanılabilecek işlevsel bir sorgu oluşturabildiniz.

Artık öğreticiyi tamamladığınıza göre, Cosmos SDK hakkındaki bilgilerinizi geliştirmeye devam edebilir ve sunduğu birçok özelliği ve yeteneği keşfedebilirsiniz. Neler yaratabileceğinizi görmek için daha karmaşık sorgular uygulamayı veya SDK'nın diğer özelliklerini denemeyi deneyebilirsiniz.
