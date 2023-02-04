# Vue frontend

Vue 3 ile blok zinciriniz için bir web uygulaması geliştirmek üzere Ignite kullanımına ilişkin bu eğitime hoş geldiniz. Ignite, hızlı bir şekilde çalışmaya başlamak için kullanılabilecek bir dizi şablon ve jeneratör sağlayarak bir blockchain uygulaması oluşturma sürecini basitleştiren bir araçtır.

[Vue 3](https://vuejs.org/) ile blockchaininiz için bir web uygulaması geliştirmek üzere Ignite kullanımına ilişkin bu eğitime hoş geldiniz. Ignite, hızlı bir şekilde çalışmaya başlamak için kullanılabilecek bir dizi şablon ve jeneratör sağlayarak bir blockchain uygulaması oluşturma sürecini basitleştiren bir araçtır.

Blok zincirinizi ve Vue şablonunuzu kurduktan sonra, bir sonraki adım bir API istemcisi oluşturmaktır. Bu, web uygulamanızdan blockchain'inizle kolayca etkileşime girmenize, veri almanıza ve işlem yapmanıza olanak tanıyacaktır. Bu eğitimin sonunda, kendi blockchaininize bağlı tamamen işlevsel bir web uygulamanız olacak.

Ön Gereksinimler:

* [Node.js](https://nodejs.org/en/)
* [Keplr](https://www.keplr.app/) Chrome eklentisi

Yeni bir blockchain projesi oluşturun:

```
ignite scaffold chain example
```

Bir Vue ön uç şablonu oluşturmak için `example` dizinine gidin ve aşağıdaki komutu çalıştırın:

Bu, `vue`dizininde yeni bir Vue projesi oluşturacaktır. Bu proje herhangi bir blok zinciri ile kullanılabilir, ancak blockchain ile etkileşim için bir API istemcisine bağlıdır. Bir API istemcisi oluşturmak için `example`dizininde aşağıdaki komutu çalıştırın:

```
ignite generate composables
```

Bu komut iki dizin oluşturur:

* `ts-client`: Blockchain'inizle etkileşim kurmak için kullanılabilecek, çerçeveden bağımsız bir TypeScript istemcisi. TypeScript istemci eğitiminde bu istemcinin nasıl kullanılacağı hakkında daha fazla bilgi edinebilirsiniz.
* `vue/src/composables`: TypeScript istemcisini saran ve Vue uygulamanızdan blockchain ile etkileşimi kolaylaştıran bir Vue 3 [bileşikleri ](https://vuejs.org/guide/reusability/composables.html)koleksiyonu.

Keplr cüzdan uzantısı yüklü olarak tarayıcınızı açın. Yeni bir hesap oluşturmak veya mevcut bir hesabı kullanmak için [talimatları](https://keplr.crunch.help/en/getting-started/creating-a-new-keplr-account) izleyin. Bir sonraki adımda ihtiyaç duyacağınız için anımsatıcı ifadeyi kaydettiğinizden emin olun.

Önem verdiğiniz varlıkların bulunduğu bir hesapla ilişkili bir anımsatıcı ifade kullanmayın. Bunu yaparsanız, bu varlıkları kaybetme riskiyle karşı karşıya kalırsınız. Geliştirme amacıyla yeni bir hesap oluşturmak iyi bir uygulamadır.&#x20;

Keplr'da kullandığınız hesabı blockchain'inizin `config.yml` dosyasına ekleyin:

```
accounts:
  - name: alice
    coins: [20000token, 200000000stake]
  - name: bob
    coins: [10000token, 100000000stake]
  - name: frank
    coins: [10000token, 100000000stake]
    mnemonic: struggle since inmate safe logic kite tag web win stay security wonder
```

`struggle since...` anımsatıcısını bir önceki adımda kaydettiğiniz anımsatıcı ile değiştirin.

Yapılandırma dosyasına bir anımsatıcı ile bir hesap eklemek, Ignite CLI'ya başlattığınızda hesabı blockchain'e eklemesini söyleyecektir. Bu, geliştirme amaçları için kullanışlıdır, ancak bunu üretimde yapmamalısınız.

Blockchain'inizi başlatmak için `example` dizinde aşağıdaki komutu çalıştırın:

Vue uygulamanızı başlatmak için `vue` dizinine gidin ve ayrı bir terminal penceresinde aşağıdaki komutu çalıştırın:

```
npm install && npm run dev
```

Tüm bağımlılıkların yüklendiğinden emin olmak için uygulamanızı `npm run dev` ile başlatmadan önce `npm install`'ı çalıştırmanız önerilir (API istemcisinin sahip oldukları da dahil, bkz. `vue/postinstall.js`).

Tarayıcınızı açın ve şu adrese gidin: [http://localhost:5173/](http://localhost:5173/).

"Cüzdanı bağla" düğmesine basın, şifrenizi Keplr'a girin ve blockchaininizi Keplr'a eklemek için "Onayla" düğmesine basın.

Keplr'ın blockchain açılır menüsünde geliştirme amacıyla kullandığınız hesabı ve "Örnek Ağ "ı seçtiğinizden emin olun. Vue uygulamanızda varlıkların bir listesini görmelisiniz.

Tebrikler! İstemci tarafında bir Vue uygulamasını başarıyla oluşturdunuz ve blockchaininize bağladınız. Projenizin geri kalanını oluşturmak için Vue uygulamanızın kaynak kodunu değiştirebilirsiniz.

Vue uygulamasının bir Cosmos zinciriyle düzgün bir şekilde etkileşime girebilmesi için doğru adres önekinin ayarlanması gerekir. Adres öneki, uygulamanın bağlı olduğu zinciri tanımlamak için kullanılır ve zincir tarafından kullanılan önekle eşleşmelidir.

Ignite varsayılan olarak `cosmos` önekiyle bir zincir oluşturur. Zincirinizi `ignite scaffold chain ... --adddress-prefix foo` ile oluşturduysanız veya zincirin kaynak kodundaki öneki manuel olarak değiştirdiyseniz, öneki Vue uygulamasında ayarlamanız gerekir.

Bir Vue uygulamasında adres önekini (prefix) ayarlamanın iki yolu vardır.

#### Ortam değişkeni kullanma[​](broken-reference) <a href="#using-an-environment-variable" id="using-an-environment-variable"></a>

`VITE_ADDRESS_PREFIX` ortam değişkenini zinciriniz için doğru adres önekine ayarlayabilirsiniz. Bu, uygulama tarafından kullanılan varsayılan öneki geçersiz kılacaktır.

`VITE_ADDRESS_PREFIX` ortam değişkenini ayarlamak için aşağıdaki komutu kullanabilirsiniz:

```
export VITE_ADDRESS_PREFIX=your-prefix
```

`your-prefix` yerine zincirinizin gerçek adres önekini yazın.

#### Kodda adres önekini ayarlama[​](broken-reference) <a href="#setting-address-prefix-in-the-code" id="setting-address-prefix-in-the-code"></a>

Alternatif olarak, `./vue/src/env.ts` dosyasındaki `prefix` değişkeninin geri dönüş değerini değiştirerek doğru adres önekini manuel olarak ayarlayabilirsiniz.

Bunu yapmak için `./vue/src/env.ts` dosyasını açın ve aşağıdaki satırı bulun:

./vue/src/env.ts

```
const prefix = process.env.VITE_ADDRESS_PREFIX || 'your-prefix';
```

`your-prefix` öğesini zincirinizin gerçek adres önekiyle değiştirin.

Dosyayı kaydedin ve değişiklikleri uygulamak için Vue uygulamasını yeniden başlatın.
