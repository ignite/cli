# Bir zincirde hata ayıklama

Ignite chain debug komutu, geliştirme sırasında sorunları bulmanıza yardımcı olabilir. İşlemin yürütülmesini kontrol ederek, değişkenleri değerlendirerek ve iş parçacığı / goroutine durumu, CPU kayıt durumu ve daha fazlası hakkında bilgi sağlayarak blok zinciri uygulamanızla etkileşime girmenizi sağlayan [Delve ](https://github.com/go-delve/delve)hata ayıklayıcısını kullanır.

### Hata Ayıklama Komutu

Hata ayıklama komutu, blok zinciri uygulama ikilisinin optimizasyonlar ve inlining kaldırılarak hata ayıklama desteği ile oluşturulmasını gerektirir. Bir hata ayıklama ikili dosyası varsayılan olarak `ignite chain serve` komutu tarafından oluşturulur veya isteğe bağlı olarak `ignite chain init` veya `ignite chain build` alt komutları çalıştırılırken `--debug` bayrağı kullanılarak oluşturulabilir.

Terminalde bir hata ayıklama oturumu başlatmak için çalıştırın:

```
ignite chain debug
```

Komut, blockchain uygulamanızı arka planda çalıştırır, ona bağlanır ve bir terminal hata ayıklayıcı kabuğu başlatır:

```
Type 'help' for list of commands.
(dlv)
```

Bu noktada blockchain uygulaması yürütmeyi engeller, böylece yürütmeye devam etmeden önce bir veya daha fazla kesme noktası ayarlayabilirsiniz.

Örneğin `<filename>:<line>` notasyonunu kullanarak istediğiniz sayıda kesme noktası ayarlamak için [break ](https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md#break)(diğer adı b) komutunu kullanın:

```
(dlv) break x/hello/client/cli/query_say_hello.go:14
```

Bu komut `x/hello/client/cli/query_say_hello.go` dosyasına 14. satırda bir kesme noktası ekler.

Tüm kesme noktaları ayarlandıktan sonra [continue ](https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md#continue)(alias c) komutunu kullanarak blok zinciri yürütmeye devam edin:

```
(dlv) continue
```

Hata ayıklayıcı kabuğu başlatacak ve bir kesme noktası tetiklendiğinde blockchain yürütmesini tekrar durduracaktır.

Hata ayıklayıcı kabuğunda, blok zinciri uygulamasını durdurmak ve hata ayıklayıcıdan çıkmak için `quit` (diğer adı `q`) veya `exit` komutlarını kullanın.

### Hata Ayıklama Sunucusu

Varsayılan terminal istemcisinin istenmediği durumlarda isteğe bağlı olarak bir hata ayıklama sunucusu başlatılabilir. Sunucu başladığında ilk olarak blockchain uygulamasını çalıştırır, ona bağlanır ve son olarak bir istemci bağlantısı bekler. Varsayılan sunucu adresi _tcp://127.0.0.1:30500_'dür ve hem JSON-RPC hem de DAP istemci bağlantılarını kabul eder.

Bir hata ayıklama sunucusu başlatmak için aşağıdaki bayrağı kullanın:

```
ignite chain debug --server
```

Bir hata ayıklama sunucusunu özel bir adresle başlatmak için aşağıdaki bayrakları kullanın:

```
ignite chain debug --server --server-address 127.0.0.1:30500
```

İstemci bağlantısı kapatıldığında hata ayıklama sunucusu otomatik olarak durur.

### İstemcilerde Hata Ayıklama

#### Gdlv: Çoklu Platform Delve Kullanıcı Arayüzü

Gdlv, Linux, Windows ve macOS için Delve'e grafiksel bir ön uçtur.

Herhangi bir yapılandırma gerektirmediği için hata ayıklama istemcisi olarak kullanmak kolaydır. Hata ayıklama sunucusu çalıştıktan ve istemci isteklerini dinledikten sonra çalıştırarak ona bağlanın:

```
gdlv connect 127.0.0.1:30500
```

Kesme noktalarının ayarlanması ve yürütmeye devam edilmesi, `break` ve `continue` komutları kullanılarak Delve ile aynı şekilde yapılır.

### Visual Studio Kodu

[Visual Studio Code](https://code.visualstudio.com/) 'u hata ayıklama istemcisi olarak kullanmak, hata ayıklama sunucusuna bağlanmasına izin vermek için bir başlangıç yapılandırması gerektirir.

[Go ](https://code.visualstudio.com/docs/languages/go)uzantısının yüklü olduğundan emin olun.

VS Code hata ayıklama, genellikle çalışma alanınızdaki `.vscode` klasörünün içinde bulunan `launch.json` dosyası kullanılarak yapılandırılır.

VS Code'u hata ayıklama istemcisi olarak ayarlamak için aşağıdaki başlatma yapılandırmasını kullanabilirsiniz:

launch.json

```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to Debug Server",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "${workspaceFolder}",
            "port": 30500,
            "host": "127.0.0.1"
        }
    ]
}
```

Alternatif olarak, "Çalıştır ve Hata Ayıkla" panelinden özel bir `launch.json` dosyası oluşturmak mümkündür. İstendiğinde "Go: Sunucuya Bağlan" etiketli Go hata ayıklayıcı seçeneğini seçin ve hata ayıklama ana bilgisayar adresini ve ardından bağlantı noktası numarasını girin.

### Örnek: Bir Blockchain Uygulamasında Hata Ayıklama

Bu kısa örnekte, yeni bir blok zinciri ve sorgu çağrıldığında bir hata ayıklama kesme noktasını tetikleyebilmek için bir sorgu oluşturmak üzere Ignite CLI kullanacağız.

Yeni bir blok zinciri oluşturun:

```
ignite scaffold chain hello
```

`hello` dizininde yeni bir sorgu scaffold'u oluşturun:

```
ignite scaffold query say-hello name --response name
```

Bir sonraki adım blockchain'in veri dizinini başlatır ve bir hata ayıklama binary'si oluşturur:

```
ignite chain init --debug
```

Başlatma işlemi tamamlandığında hata ayıklayıcı kabuğunu başlatın:

```
ignite chain debug
```

Hata ayıklayıcı kabuğunda, `SayHello` işlevi çağrıldığında tetiklenecek bir kesme noktası oluşturun ve ardından yürütmeye devam edin:

```
(dlv) break x/hello/keeper/query_say_hello.go:12
(dlv) continue
```

Farklı bir terminalden sorguyu çağırmak için `hellod` ikili dosyasını kullanın:

```
hellod query hello say-hello bob
```

Kesme noktası tetiklendiğinde bir hata ayıklayıcı kabuğu başlatılacaktır:

```
     7:     "google.golang.org/grpc/codes"
     8:     "google.golang.org/grpc/status"
     9:     "hello/x/hello/types"
    10: )
    11:
=>  12: func (k Keeper) SayHello(goCtx context.Context, req *types.QuerySayHelloRequest) (*types.QuerySayHelloResponse, error) {
    13:     if req == nil {
    14:         return nil, status.Error(codes.InvalidArgument, "invalid request")
    15:     }
    16:
    17:     ctx := sdk.UnwrapSDKContext(goCtx)
```

Bundan sonra, yürütmeyi kontrol etmek ve değerleri yazdırmak için `next` (takma ad `n`) veya `print` (takma ad `p`) gibi Delve komutlarını kullanabilirsiniz. Örneğin, _name_ argüman değerini yazdırmak için print komutunu ve ardından "`req.Name`" komutunu kullanın:

```
(dlv) print req.Name
"bob"
```

Son olarak, blockchain uygulamasını durdurmak ve hata ayıklama oturumunu bitirmek için `quit` (diğer adı `q`) kullanın.
