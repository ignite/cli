# React frontend

React ile blockchaininiz için bir web uygulaması geliştirmek üzere Ignite kullanımına ilişkin bu eğitime hoş geldiniz. Ignite, hızlı bir şekilde çalışmaya başlamak için kullanılabilecek bir dizi şablon ve jeneratör sağlayarak bir blockchain uygulaması oluşturma sürecini basitleştiren bir araçtır.

Ignite'ın özelliklerinden biri, kullanıcı arayüzleri oluşturmak için popüler bir JavaScript çerçevesi olan [React](https://reactjs.org/)'i desteklemesidir. Bu eğitimde, yeni bir blockchain oluşturmak ve bir React ön uç şablonunu iskelelemek için Ignite'ı nasıl kullanacağınızı öğreneceksiniz. Bu size web uygulamanız için temel bir temel sağlayacak ve uygulamanızın geri kalanını oluşturmaya başlamanızı kolaylaştıracaktır.

Blok zincirinizi ve React şablonunuzu kurduktan sonra, bir sonraki adım bir API istemcisi oluşturmaktır. Bu, web uygulamanızdan blockchain'inizle kolayca etkileşime girmenizi sağlayarak veri almanıza ve işlem yapmanıza olanak tanıyacaktır. Bu eğitimin sonunda, kendi blok zincirinize bağlı tamamen işlevsel bir web uygulamasına sahip olacaksınız.

Ön Gereksinimler:

* [Node.js](https://nodejs.org/en/)
* [Keplr](https://www.keplr.app/) Chrome eklentisi

Yeni bir blockchain projesi oluşturun:

```
ignite scaffold chain example
```

To create a React frontend template, go to the `example` directory and run the following command:

Bu, `react` dizininde yeni bir React projesi oluşturacaktır. Bu proje herhangi bir blockchain ile kullanılabilir, ancak blockchain ile etkileşim için bir API istemcisine bağlıdır. Bir API istemcisi oluşturmak için `example` dizininde aşağıdaki komutu çalıştırın:

Bu komut iki dizin oluşturur:

* `ts-client`: blok zincirinizle etkileşim kurmak için kullanılabilecek, çerçeveden bağımsız bir TypeScript istemcisi. TypeScript istemci eğitiminde bu istemcinin nasıl kullanılacağı hakkında daha fazla bilgi edinebilirsiniz.
* `react/src/hooks`: TypeScript istemcisini saran ve React uygulamanızdan blok zincirinizle etkileşimi kolaylaştıran bir [React Hooks](https://reactjs.org/docs/hooks-intro.html) koleksiyonu.

Keplr cüzdan uzantısı yüklü olarak tarayıcınızı açın. Yeni bir hesap oluşturmak veya mevcut bir hesabı kullanmak için [talimatları](https://keplr.crunch.help/en/getting-started/creating-a-new-keplr-account) izleyin. Bir sonraki adımda ihtiyaç duyacağınız için anımsatıcı ifadeyi kaydettiğinizden emin olun.

Önem verdiğiniz varlıkların bulunduğu bir hesapla ilişkili bir anımsatıcı ifade kullanmayın. Bunu yaparsanız, bu varlıkları kaybetme riskiyle karşı karşıya kalırsınız. Geliştirme amacıyla yeni bir hesap oluşturmak iyi bir uygulamadır.

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

Blockchaini başlatmak için `example` dizininde aşağıdaki komutu çalıştırın:&#x20;

React uygulamanızı başlatmak için `react` dizinine gidin ve ayrı bir terminal penceresinde aşağıdaki komutu çalıştırın:

```
npm install && npm run dev
```

Tüm bağımlılıkların yüklendiğinden emin olmak için uygulamanızı `npm run dev` ile başlatmadan önce `npm install`'ı çalıştırmanız önerilir (API istemcisinin sahip oldukları dahil, bkz. `react/postinstall.js)`.

Tarayıcınızı açın ve [http://localhost:5173/](http://localhost:5173/) adresine gidin.

"Cüzdanı bağla" düğmesine basın, şifrenizi Keplr'a girin ve blockchaininizi Keplr'a eklemek için "Onayla" düğmesine basın.

Keplr'ın blok zinciri açılır menüsünde geliştirme amacıyla kullandığınız hesabı ve "Örnek Ağ "ı seçtiğinizden emin olun. React uygulamanızda varlıkların bir listesini görmelisiniz.

Tebrikler! İstemci tarafında bir React uygulamasını başarıyla oluşturdunuz ve blockchaininize bağladınız. Projenizin geri kalanını oluşturmak için React uygulamanızın kaynak kodunu değiştirebilirsiniz.
