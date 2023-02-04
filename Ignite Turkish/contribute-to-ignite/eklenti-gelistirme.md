# Eklenti Geliştirme

Bir eklenti oluşturmak ve projenizde hemen kullanmak kolaydır. Önce projenizin dışında bir dizin seçin ve :

```
$ ignite plugin scaffold my-plugin
```

Bu, eklentinin kodunu içeren yeni bir `my-plugin` dizini oluşturacak ve eklentinizin `ignite` komutuyla nasıl kullanılacağına ilişkin bazı talimatlar verecektir. Aslında, bir eklenti yolu yerel bir dizin olabilir ve bunun çeşitli faydaları vardır:

* eklentinizin geliştirilmesi sırasında bir git deposu kullanmanıza gerek yoktur.
* Kaynak dosyalar eklenti ikilisinden daha eskiyse, projenizde `ignite` ikilisini her çalıştırdığınızda eklenti yeniden derlenir.

Böylece eklenti geliştirme iş akışı bu kadar basittir :

1. `ignite plugin scaffold my-plugin` ile bir eklentiyi iskeleleyin&#x20;
2. `ignite plugin add -g /path/to/my-plugin` aracılığıyla yapılandırmanıza ekleyin&#x20;
3. eklenti kodunu güncelle&#x20;
4. eklentiyi derlemek ve çalıştırmak için `ignite my-plugin` ikili dosyasını çalıştırın.&#x20;
5. 3'e geri dönün.

Eklentiniz hazır olduğunda, onu bir git deposunda yayınlayabilirsiniz ve topluluk `ignite plugin add github.com/foo/my-plugin`'i çağırarak onu kullanabilir.

Şimdi eklentinizin kodunu nasıl güncelleyeceğinizi detaylandıralım.

`ignite` eklenti sistemi `github.com/hashicorp/go-` eklentisi, önceden tanımlanmış bir arayüzü uygulamak anlamına gelir:

ignite/services/plugin/interface.go

```
// An ignite plugin must implements the Plugin interface.
type Interface interface {
    // Manifest declares the plugin's Command(s) and Hook(s).
    Manifest() (Manifest, error)

    // Execute will be invoked by ignite when a plugin Command is executed.
    // It is global for all commands declared in Manifest, if you have declared
    // multiple commands, use cmd.Path to distinguish them.
    Execute(cmd ExecutedCommand) error

    // ExecuteHookPre is invoked by ignite when a command specified by the Hook
    // path is invoked.
    // It is global for all hooks declared in Manifest, if you have declared
    // multiple hooks, use hook.Name to distinguish them.
    ExecuteHookPre(hook ExecutedHook) error

    // ExecuteHookPost is invoked by ignite when a command specified by the hook
    // path is invoked.
    // It is global for all hooks declared in Manifest, if you have declared
    // multiple hooks, use hook.Name to distinguish them.
    ExecuteHookPost(hook ExecutedHook) error

    // ExecuteHookCleanUp is invoked by ignite when a command specified by the
    // hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
    // execution status of the command and hooks.
    // It is global for all hooks declared in Manifest, if you have declared
    // multiple hooks, use hook.Name to distinguish them.
    ExecuteHookCleanUp(hook ExecutedHook) error
}
```

İskeletlenen kod zaten bu arayüzü uygular, sadece yöntemlerin gövdesini güncellemeniz gerekir.

İşte `Manifest` yapısı :

ignite/services/plugin/interface.go

```
type Manifest struct {
    Name string
    // Commands contains the commands that will be added to the list of ignite
    // commands. Each commands are independent, for nested commands use the
    // inner Commands field.
    Commands []Command
    // Hooks contains the hooks that will be attached to the existing ignite
    // commands.
    Hooks []Hook
    // SharedHost enables sharing a single plugin server across all running instances
    // of a plugin. Useful if a plugin adds or extends long running commands
    //
    // Example: if a plugin defines a hook on `ignite chain serve`, a plugin server is instanciated
    // when the command is run. Now if you want to interact with that instance from commands
    // defined in that plugin, you need to enable `SharedHost`, or else the commands will just
    // instantiate separate plugin servers.
    //
    // When enabled, all plugins of the same `Path` loaded from the same configuration will
    // attach it's rpc client to a an existing rpc server.
    //
    // If a plugin instance has no other running plugin servers, it will create one and it will be the host.
    SharedHost bool `yaml:"shared_host"`
}
```

Eklentinizin kodunda, `Manifest` yöntemi zaten örnek olarak önceden tanımlanmış bir ~~Manifest~~ yapısı döndürür. İhtiyacınıza göre uyarlayın.

Eklentiniz`ignite`'a bir veya daha fazla yeni komut eklerse, `Commands` alanını besler.

Eklentiniz mevcut komutlara özellikler eklerse, `Hooks` alanını besler.

Elbette bir eklenti `Commands` ve `Hooks` bildirebilir.

Bir eklenti `SharedHost`'u `true` olarak ayarlayarak bir ana bilgisayar sürecini de paylaşabilir. `SharedHost`, bir eklenti uzun süre çalışan komutlara bağlanıyorsa veya bunları bildiriyorsa tercih edilir. Aynı eklenti bağlamından çalıştırılan komutlar aynı eklenti sunucusuyla etkileşime girer. Yürütülen tüm komutların aynı sunucu örneğini paylaşmasına izin vererek paylaşılan yürütme bağlamı sağlar.

**Yeni komut ekleme**

Eklenti komutları, kayıtlı bir eklenti tarafından ignite cli'ye eklenen özel komutlardır. Komutlar, ignite tarafından önceden tanımlanmamış herhangi bir yolda olabilir. Tüm eklenti komutları `ignite` komut kökünü genişletecektir.

Örneğin, eklentinizin `ignite scaffold`'a yeni bir `oracle` komutu eklediğini varsayalım, `Manifest()` yöntemi şöyle görünecektir :

```
func (p) Manifest() (plugin.Manifest, error) {
    return plugin.Manifest{
        Name: "oracle",
        Commands: []plugin.Command{
            {
                Use:   "oracle [name]",
                Short: "Scaffold an oracle module",
                Long:  "Long description goes here...",
                // Optionnal flags is required
                Flags: []plugin.Flag{
                    {Name: "source", Type: plugin.FlagTypeString, Usage: "the oracle source"},
                },
                // Attach the command to `scaffold`
                PlaceCommandUnder: "ignite scaffold",
            },
        },
    }, nil
}
```

Eklenti yürütmesini güncellemek için, eklenti `Execute` komutunu değiştirmeniz gerekir, örneğin :

```
func (p) Execute(cmd plugin.ExecutedCommand) error {
    if len(cmd.Args) == 0 {
        return fmt.Errorf("oracle name missing")
    }
    var (
        name      = cmd.Args[0]
        source, _ = cmd.Flags().GetString("source")
    )
    // Read chain information
    c, err := getChain(cmd)
    if err != nil {
        return err
    }

    //...
}
```

Ardından, eklentiyi çalıştırmak için `ignite scaffold oracle`'ı çalıştırın.

Eklenti `Hooks`'ları, mevcut ignite komutlarının yeni işlevlerle genişletilmesine olanak tanır. Kancalar, bir komut çalıştırıldıktan sonra veya önce özel komut dosyaları çalıştırmaya gerek kalmadan işlevselliği düzene sokmak istediğinizde kullanışlıdır. bu, bir zamanlar hataya açık olan veya hep birlikte unutulan süreçleri düzene sokabilir.

Aşağıda, kayıtlı bir `ignite` komutu üzerinde çalışacak kancalar tanımlanmıştır

| İsim     | Açıklama                                                                                                                                    |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| Pre      | `PreRun` kapsamında bir komutun ana işlevi çağrılmadan önce çalışır                                                                         |
| Post     | `PostRun` kapsamında bir komutun ana işlevi çağrıldıktan sonra çalışır                                                                      |
| Clean Up | Bir komutun ana işlevi çağrıldıktan sonra çalışır. komut bir hata döndürürse, yürütmeyi garanti etmek için hata döndürülmeden önce çalışır. |

Not: Bir hook ön adımda bir hataya neden olursa komut çalışmaz, bu da `post` ve `clean up` işlemlerinin yürütülmemesine neden olur.

Aşağıda bir `hook` tanımı örneği verilmiştir.

```
func (p) Manifest() (plugin.Manifest, error) {
    return plugin.Manifest{
        Name: "oracle",
        Hooks: []plugin.Hook{
            {
                Name:        "my-hook",
                PlaceHookOn: "ignite chain build",
            },
        },
    }, nil
}

func (p) ExecuteHookPre(hook plugin.ExecutedHook) error {
    switch hook.Name {
    case "my-hook":
        fmt.Println("I'm executed before ignite chain build")
    default:
        return fmt.Errorf("hook not defined")
    }
    return nil
}

func (p) ExecuteHookPost(hook plugin.ExecutedHook) error {
    switch hook.Name {
    case "my-hook":
        fmt.Println("I'm executed after ignite chain build (if no error)")
    default:
        return fmt.Errorf("hook not defined")
    }
    return nil
}

func (p) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
    switch hook.Name {
    case "my-hook":
        fmt.Println("I'm executed after ignite chain build (regardless errors)")
    default:
        return fmt.Errorf("hook not defined")
    }
    return nil
}
```

Yukarıda, bir hook'un bir `Name` ve bir `PlaceHookOn`'a sahip olduğu `Command`'a benzer bir tanım görebiliriz. `Execute*` yöntemlerinin doğrudan kancanın her bir yaşam döngüsüne eşlendiğini fark edeceksiniz. Eklenti içinde tanımlanan tüm hook'lar bu metotları çağıracaktır.
