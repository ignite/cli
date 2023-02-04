# App Design

Bu bölümde, chain'ler arası borsa modülünün nasıl tasarlandığını öğreneceksiniz. Modülde emir defterleri, alış emirleri ve satış emirleri vardır.

İlk olarak, bir token çifti için bir emir defteri oluşturun.

Bir emir defteri oluşturulduktan sonra, bu token çifti için alım ve satım emirleri oluşturabilirsiniz.

Modül, Blok Zincirler Arası İletişim protokolü [IBC](https://github.com/cosmos/ibc/blob/old/ibc/2\_IBC\_ARCHITECTURE.md)'yi kullanır. Modül, IBC'yi kullanarak birden fazla blockchain'in etkileşime girebilmesi ve tokenlarını takas edebilmesi için emir defterleri oluşturabilir.

Bir blockchain'den bir token ve başka bir blockchain'den başka bir token ile bir emir defteri çifti oluşturursunuz. Bu eğitimde, oluşturduğunuz modülü `dex` modülü olarak adlandırın.

> Bir kullanıcı `dex` modülü ile bir token takas ettiğinde, diğer blockchain üzerinde o tokena ait bir kupon alınır. Bu kupon, `ibc-transfer`inin nasıl oluşturulduğuna benzer. Bir blockchain modülü, bir blockchainin yeni tokenlarını var etme hakkına sahip olmadığından, hedef chain'deki token kilitlenir ve alıcı bu tokenın bir `vouche`'ını alır.

Bu süreç, orijinal tokenin kilidini açmak için `voucher` yakıldığında tersine çevrilebilir. Bu değişim süreci eğitim boyunca daha ayrıntılı olarak açıklanmaktadır.

### Tasarım Varsayımı

Herhangi bir chain çifti arasında herhangi bir token değişimi için bir emir defteri oluşturulabilir.

* Her iki blockchain de `dex` modülünün kurulu ve çalışır durumda olmasını gerektirir.
* Aynı anda bir token çifti için yalnızca bir sipariş defteri olabilir.

Belirli bir chain kendi yerel tokenından yeni coin basamaz.

Bu modül Cosmos SDK'daki [ibc transfer](https://github.com/cosmos/ibc-go/tree/main/modules/apps/transfer) modülünden esinlenmiştir. Bu eğitimde oluşturduğunuz dex modülü, kupon oluşturma gibi benzerliklere sahiptir.

Bununla birlikte, oluşturduğunuz yeni dex modülü daha karmaşıktır çünkü aşağıdakilerin oluşturulmasını destekler:

* Gönderilecek çeşitli paket türleri
* Tedavi edilecek çeşitli teşekkür türleri
* Bir paketin alındığında, zaman aşımında ve daha fazlasında nasıl ele alınacağına ilişkin daha karmaşık mantık

### Interchain Borsasına Genel Bakış

İki blockchain'iniz olduğunu varsayalım: Venüs ve Mars.

Venüs'teki yerel token `venuscoin`'dir.

Mars'taki yerel token `marscoin`'dir.

Bir token Mars'tan Venüs'e takas edildiğinde:

* Venüs blockchaini'nde `ibc/B5CB286...A7B21307F` gibi görünen bir denoma sahip bir IBC `voucher` jetonu vardır.
* `ibc/`'den sonraki uzun karakter dizisi, IBC kullanılarak aktarılan bir tokenin denom izleme karmasıdır.

Blockchain'in API'sini kullanarak bu hash'ten bir denom izi elde edebilirsiniz. Denom izi bir `base_denom` ve bir `path`'den oluşur. Bizim örneğimizde:

* `Base_denom` `marscoin`'dir.
* `Path`, token'ın aktarıldığı port ve kanal çiftlerini içerir.

Tek atlamalı bir aktarım için `path`, `transfer/channel-0` ile tanımlanır.

[ICS 20 Fungible Token Transfer](https://github.com/cosmos/ibc/tree/main/spec/app/ics-020-fungible-token-transfer)'de token yolları hakkında daha fazla bilgi edinin.

**Not**: Bu token `ibc/Venus/marscoin` aynı emir defteri kullanılarak geri satılamaz. Değişimi "tersine çevirmek" ve Mars tokenini geri almak istiyorsanız `ibc/Venus/marscoin`'den `marscoin`'e transfer için yeni bir emir defteri oluşturmanız ve kullanmanız gerekir.

### Sipariş Defterlerinin Tasarımı

Tipik bir borsa olarak, yeni bir çift, `marscoin` satma emirleri veya `venuscoin` satın alma emirleri içeren bir emir defterinin oluşturulması anlamına gelir. Burada, iki zinciriniz var ve bu veri yapısı Mars ve Venüs arasında bölünmelidir.

* Mars zincirindeki kullanıcılar `marscoin` satar.
* Venüs zincirindeki kullanıcılar `marscoin` satın alır.

Bu nedenle, temsil ediyoruz:

* Mars zincirinde `marscoin` satmak için tüm siparişler.
* Venüs zincirinde `marscoin` satın almak için tüm emirler.

Bu örnekte, Mars blok zinciri satış emirlerini, Venüs blok zinciri ise alış emirlerini tutar.

### Tokenları Geri Takas Etme

`ibc-transfer` gibi, her blok zinciri diğer blok zincirinde oluşturulan token voucher'inin bir izini tutar.

Mars blok zinciri Venüs zincirine `marscoin` satarsa ve `ibc/Venüs/marscoin` Venüs'te basılırsa, `ibc/Venüs/marscoin` Mars'a geri satılırsa, token kilidi açılır ve alınan token `marscoin` olur.

### Özellikler

Zincirler arası değişim modülü tarafından desteklenen özellikler şunlardır:

* İki zincir arasında bir token çifti için bir değişim emri defteri oluşturma
* Kaynak zincirinde satış emirleri gönderin
* Hedef zincire satın alma emirleri gönderin
* Satış veya alış emirlerini iptal etme
