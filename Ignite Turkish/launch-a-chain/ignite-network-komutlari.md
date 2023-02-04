# Ignite Network komutları

`ignite network`, Ignite Zinciri ile etkileşime girerek egemen Cosmos blockchain'lerinin başlatılmasını koordine etmeyi sağlar.

Bir Cosmos blockchain'i başlatmak için birinin koordinatör ve diğerlerinin de validatör olması gerekir. Bunlar sadece rollerdir, herkes koordinatör veya doğrulayıcı olabilir.

* Bir koordinatör, Ignite blockchain'inde başlatılacak bir zincir hakkında bilgi yayınlar, validatör taleplerini onaylar ve başlatmayı koordine eder.
* Validatörler bir zincire katılmak için talepler gönderir ve bir Blockchain başlatılmaya hazır olduğunda düğümlerini başlatır.

CLI ile başlatmak, `ignite network` ad alanını kullanarak CLI ile birkaç kısa komut kadar basit olabilir.

> **NOTE:** `ignite n` can also be used as a shortcut for `ignite network`.

Zincirinizle ilgili bilgileri bir koordinatör olarak yayınlamak için aşağıdaki komutu çalıştırın (URL, Cosmos SDK zincirinin bulunduğu bir depoyu işaret etmelidir):

```
ignite network chain publish github.com/ignite/example
```

Bu komut, aşağıdaki komutlarda kullanacağınız launch tanımlayıcısını döndürecektir. Bu tanımlayıcının 42 olduğunu varsayalım. Ardından, doğrulayıcılardan node'larını başlatmalarını ve ağa katılma talebinde bulunmalarını isteyin. Bir test ağı için CLI tarafından önerilen varsayılan değerleri kullanabilirsiniz.

```
ignite network chain init 42
ignite network chain join 42 --amount 95000000stake
```

Koordinatör olarak tüm validatör taleplerini listeleyin:

```
ignite network request list 42
```

Validatör taleplerini onaylayın:

```
ignite network request approve 42 1,2
```

Validatör setinde ihtiyacınız olan tüm validatörleri onayladıktan sonra zincirin başlatılmaya hazır olduğunu duyurun:

```
ignite network chain launch 42
```

Validatörler artık node'larını launch için hazırlayabilirler:

```
ignite network chain prepare 42
```

Bu komutun çıktısı, bir validatörün node'unu başlatmak için kullanacağı bir komut gösterecektir, örneğin `exampled --home ~/.example`. Yeterli sayıda validatör node'larını başlattıktan sonra, bir Blockchain canlı olacaktır.

***

Sonraki iki bölüm, bir koordinatörün bir zincir launch'ını koordine etme ve bir validatör olarak bir zincir launch'ına katılma süreci hakkında daha fazla bilgi sağlamaktadır.
