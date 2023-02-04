# Başlangıç

Bu eğitimde, yeni bir blok zinciri oluşturmak için Ignite CLI kullanacağız. Ignite CLI, kullanıcıların hızlı ve kolay bir şekilde blok zinciri ağları oluşturmasına olanak tanıyan bir komut satırı arayüzüdür. Ignite CLI kullanarak, gerekli tüm bileşenleri manuel olarak ayarlamak zorunda kalmadan hızlı bir şekilde yeni bir blok zinciri oluşturabiliriz.

Blok zincirimizi Ignite CLI ile oluşturduktan sonra, oluşturulan dizin yapısına ve dosyalara bir göz atacağız. Bu bize blok zincirinin nasıl organize edildiğini ve blok zincirinin farklı bileşenlerinin birbirleriyle nasıl etkileşime girdiğini anlamamızı sağlayacaktır.

Bu eğitimin sonunda, yeni bir blok zinciri oluşturmak için Ignite CLI'nın nasıl kullanılacağına dair temel bir anlayışa sahip olacak ve bir blok zincirini oluşturan dizin yapısı ve dosyalar hakkında üst düzey bir anlayışa sahip olacaksınız. Bu bilgi, blok zinciri geliştirme dünyasını keşfetmeye devam ederken faydalı olacaktır.

### Yeni bir blok zinciri oluşturma

Ignite ile yeni bir blok zinciri projesi oluşturmak için aşağıdaki komutu çalıştırmanız gerekecektir:

```
ignite scaffold chain example
```

[ignite scaffold chain](https://docs.ignite.com/references/cli#ignite-scaffold-chain) komutu yeni bir dizin örneğinde yeni bir blok zinciri oluşturacaktır.

Yeni blok zinciri Cosmos SDK çerçevesi kullanılarak oluşturulur ve bir dizi işlevsellik sağlamak için birkaç standart modülü içe aktarır. Bu modüller arasında delege edilmiş bir Proof-of-Stake konsensüs mekanizması sağlayan `staking`, hesaplar arasında değiştirilebilir token transferlerini kolaylaştırmak için banka ve zincir üzerinde yönetişim için gov bulunmaktadır. Bu modüllere ek olarak, blok zinciri Cosmos SDK çerçevesinden diğer modülleri de içe aktarır.

`example` dizin, Cosmos SDK blok zincirinin yapısını oluşturan oluşturulmuş dosya ve dizinleri içerir. Bu dizin, diğerlerinin yanı sıra zincirin yapılandırması, uygulama mantığı ve testler için dosyalar içerir. Geliştiricilerin hızlı bir şekilde yeni bir Cosmos SDK blok zinciri kurmaları ve istedikleri işlevselliği bunun üzerine inşa etmeleri için bir başlangıç noktası sağlar.

Varsayılan olarak Ignite, `x/` dizininde oluşturulmakta olan blok zinciri (bu durumda `example`) ile aynı ada sahip yeni bir boş özel modül oluşturur. Bu modül kendi başına herhangi bir işlevselliğe sahip değildir, ancak uygulamanızın özelliklerini oluşturmak için bir başlangıç noktası olarak hizmet edebilir. Bu modülü oluşturmak istemiyorsanız, atlamak için `--no-module` bayrağını kullanabilirsiniz.

### Dizin yapısı

Ignite CLI'nin projeniz için ne oluşturduğunu anlamak için `example/` dizininin içeriğini inceleyebilirsiniz.

`app/` dizini blok zincirinin farklı parçalarını birbirine bağlayan dosyaları içerir. Bu dizindeki en önemli dosya, blok zincirinin tip tanımını ve onu oluşturma ve başlatma işlevlerini içeren `app.go` dosyasıdır. Bu dosya, blok zincirinin çeşitli bileşenlerini birbirine bağlamaktan ve birbirleriyle nasıl etkileşime gireceklerini tanımlamaktan sorumludur.

`cmd/` dizini, derlenmiş ikilinin komut satırı arayüzünden (CLI) sorumlu ana paketi içerir. Bu paket, CLI'dan çalıştırılabilecek komutları ve bunların nasıl yürütülmesi gerektiğini tanımlar. Geliştiricilerin ve kullanıcıların blok zinciri ile etkileşime girmesi ve blok zinciri durumunu sorgulamak veya işlem göndermek gibi çeşitli görevleri yerine getirmesi için bir yol sağladığından blok zinciri projesinin önemli bir parçasıdır.

`docs/` dizini proje belgelerini saklamak için kullanılır. Varsayılan olarak bu dizin, bir yazılım projesinin API'sini tanımlamak için makine tarafından okunabilir bir format olan bir OpenAPI belirtim dosyası içerir. OpenAPI belirtimi, proje için otomatik olarak insan tarafından okunabilir belgeler oluşturmak için kullanılabileceği gibi, diğer araç ve hizmetlerin API ile etkileşime girmesi için bir yol da sağlayabilir. `docs/` dizini, projeyle ilgili tüm ek belgeleri saklamak için kullanılabilir.

`proto/` dizini, blok zincirinin veri yapısını tanımlamak için kullanılan protokol buffer dosyalarını içerir. Protokol buffer'ları, yapılandırılmış verilerin serileştirilmesi için dil ve platformdan bağımsız bir mekanizmadır ve genellikle blok zinciri ağları gibi dağıtılmış sistemlerin geliştirilmesinde kullanılır. `Proto/` dizinindeki protokol buffer dosyaları, blok zinciri tarafından kullanılan veri yapılarını ve mesajları tanımlar ve blok zinciri ile etkileşimde bulunmak için kullanılabilecek çeşitli programlama dilleri için kod üretmek için kullanılır. Cosmos SDK bağlamında, protokol buffer dosyaları, blok zinciri tarafından gönderilip alınabilecek belirli veri türlerinin yanı sıra blok zincirinin işlevselliğine erişmek için kullanılabilecek belirli RPC uç noktalarını tanımlamak için kullanılır.

`testutil/` dizini test için kullanılan yardımcı fonksiyonları içerir. Bu fonksiyonlar, blok zinciri için testler yazarken ihtiyaç duyulan test hesapları oluşturma, işlem oluşturma ve blok zincirinin durumunu kontrol etme gibi yaygın görevleri gerçekleştirmek için uygun bir yol sağlar. Geliştiriciler `testutil/` dizinindeki yardımcı fonksiyonları kullanarak testleri daha hızlı ve verimli bir şekilde yazabilir ve testlerinin kapsamlı ve etkili olmasını sağlayabilirler.

`x/` dizini, blok zincirine eklenen özel Cosmos SDK modüllerini içerir. Standart Cosmos SDK modülleri, Cosmos SDK tabanlı blok zincirleri için stake etme ve yönetişim desteği gibi ortak işlevler sağlayan önceden oluşturulmuş bileşenlerdir. Özel modüller ise blok zinciri projesi için özel olarak geliştirilen ve projeye özgü işlevsellik sağlayan modüllerdir.

`config.yml` dosyası, geliştirme sırasında blok zincirini özelleştirmek için kullanılabilecek bir yapılandırma dosyasıdır. Bu dosya, ağın kimliği, hesap bakiyeleri ve node parametreleri gibi blok zincirinin çeşitli yönlerini kontrol eden ayarları içerir.

`.github` dizini, bir blok zinciri ikili dosyasını otomatik olarak oluşturmak ve yayınlamak için kullanılabilecek bir GitHub Actions iş akışı içerir. GitHub Actions, geliştiricilerin projelerini oluşturma, test etme ve dağıtma dahil olmak üzere yazılım geliştirme iş akışlarını otomatikleştirmelerini sağlayan bir araçtır. `.github` dizinindeki iş akışı, blok zinciri ikilisini oluşturma ve yayınlama sürecini otomatikleştirmek için kullanılır, bu da geliştiriciler için zaman ve emek tasarrufu sağlayabilir.

`Readme.md` dosyası, blok zinciri projesine genel bir bakış sağlayan bir benioku dosyasıdır. Bu dosya tipik olarak projenin adı ve amacı gibi bilgilerin yanı sıra blok zincirinin nasıl oluşturulacağı ve çalıştırılacağına ilişkin talimatları da içerir. Geliştiriciler ve kullanıcılar `readme.md` dosyasını okuyarak blok zinciri projesinin amacını ve yeteneklerini hızlı bir şekilde anlayabilir ve kullanmaya başlayabilirler.

### Bir blok zinciri node'u başlatma

Bir blok zinciri node'unu geliştirme modunda başlatmak için aşağıdaki komutu çalıştırabilirsiniz:

```
ignite chain serve
```

[ignite chain serve](https://docs.ignite.com/references/cli#ignite-scaffold-chain) komutu, bir blok zinciri node'unu geliştirme modunda başlatmak için kullanılır. Önce `ignite chain build` komutunu kullanarak binary'yi derler ve yükler, ardından `ignite chain init` komutunu kullanarak tek bir validatör için blok zincirinin veri dizinini başlatır. Bundan sonra, node'u yerel olarak başlatır ve otomatik kod yeniden yüklemeyi etkinleştirir, böylece koddaki değişiklikler node'u yeniden başlatmak zorunda kalmadan çalışan blok zincirine yansıtılabilir. Bu, blok zincirinin daha hızlı geliştirilmesine ve test edilmesine olanak tanır.

Tebrikler! 🥳 Ignite CLI kullanarak yepyeni bir Cosmos blok zincirini başarıyla oluşturdunuz. Bu blok zinciri, delegated proof of stake (DPoS) konsensüs algoritmasını kullanır ve token transferleri, yönetişim ve enflasyon için bir dizi standart modülle birlikte gelir. Artık Cosmos blok zinciriniz hakkında temel bir anlayışa sahip olduğunuza göre, özel işlevler oluşturmaya başlamanın zamanı geldi. Aşağıdaki eğitimlerde, özel modülleri nasıl oluşturacağınızı ve blok zincirinize yeni özellikleri nasıl ekleyeceğinizi öğrenerek benzersiz ve güçlü bir merkezi olmayan uygulama oluşturmanıza olanak tanıyacaksınız.
