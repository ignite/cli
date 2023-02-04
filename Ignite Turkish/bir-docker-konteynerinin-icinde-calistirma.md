# Bir Docker konteynerinin içinde çalıştırma

Ignite CLI ikili dosyasını doğrudan makinenize yüklemeden Ignite CLI'yı bir Docker konteynerinin içinde çalıştırabilirsiniz.

Ignite CLI'yi Docker'da çalıştırmak çeşitli nedenlerle yararlı olabilir; test ortamınızı izole etmek, Ignite CLI'yi desteklenmeyen bir işletim sisteminde çalıştırmak veya Ignite CLI'nin farklı bir sürümünü yüklemeden denemek.

Docker konteynerleri sanal makineler gibidir, çünkü içlerinde çalışan programlara yalıtılmış bir ortam sağlarlar. Bu durumda, Ignite CLI'yi yalıtılmış bir ortamda çalıştırabilirsiniz.

Deneme ve dosya sistemi etkisi Docker örneği ile sınırlıdır. Ana makine, kapsayıcıdaki değişikliklerden etkilenmez.

Docker'ın yüklenmiş olması gerekir. [Docker ile Başlarken](https://www.docker.com/get-started) bölümüne bakın.

Docker konteynerinizde bir zincir oluşturduktan ve başlattıktan sonra, tüm Ignite CLI komutları kullanılabilir. Komutları `docker run -ti ignite/cli` komutundan sonra yazmanız yeterlidir. Örneğin:

```
docker run -ti ignitehq/cli -h
docker run -ti ignitehq/cli scaffold chain planet
docker run -ti ignitehq/cli chain serve
```

Docker yüklendiğinde, tek bir komutla bir blockchain oluşturabilirsiniz.

Ignite CLI ve Ignite CLI ile hizmet verdiğiniz zincirler bazı dosyaları saklar. CLI ikilisini doğrudan kullanırken, bu dosyalar `$HOME/.ignite` ve `$HOME/.cache` içinde bulunur, ancak Docker bağlamında `$HOME`'dan farklı bir dizin kullanmak daha iyidir, bu yüzden `$HOME/sdh` kullanıyoruz. Bu klasör aşağıdaki docker komutlarından önce manuel olarak oluşturulmalıdır, aksi takdirde Docker bunu root kullanıcısı ile oluşturur.

Konteynerdeki `/apps` dizininde bir blockchain `planet`'i iskelesi oluşturmak için bu komutu bir terminal penceresinde çalıştırın:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps ignitehq/cli:0.25.2 scaffold chain planet
```

Sabırlı olun, bu komutun çalışması bir iki dakika sürer çünkü her şeyi sizin için yapar:

* `ignitehq/cli:0.25.2` imajından çalışan bir konteyner oluşturur.
* İmaj içindeki Ignite CLI binary'sini çalıştırır.
* `-v $HOME/sdh:/home/tendermint` yerel bilgisayarınızdaki (ana makine) `$HOME/sdh` dizinini konteyner içindeki `/home/tendermint` ev dizinine eşler.
*   `-v $PWD:/apps,` ana makinedeki terminal penceresindeki geçerli dizini konteynerdeki `/apps` diziniyle eşler. İsteğe bağlı olarak `$PWD` yerine mutlak bir yol belirtebilirsiniz.

    `w` ve `-v` birlikte kullanıldığında ana makinede dosya kalıcılığı sağlanır. Docker konteynerindeki uygulama kaynak kodu ana makinenin dosya sistemine yansıtılır.

    Not: `-w` ve `-v` bayrakları için dizin adı `/app` dışında bir ad olabilir, ancak her iki bayrak için de aynı dizin belirtilmelidir. `w` ve `-v`'yi atlarsanız, değişiklikler yalnızca kapsayıcıda yapılır ve kapsayıcı kapatıldığında kaybolur.

Blockchain node'unu yeni oluşturduğunuz Docker konteynerinde başlatmak için şu komutu çalıştırın:

```
docker run -ti -v $HOME/sdh:/home/tendermint -v $PWD:/apps -p 1317:1317 -p 26657:26657 ignitehq/cli:0.25.2 chain serve -p planet
```

Bu komut aşağıdakileri yapar:

* `-v $HOME/sdh:/home/tendermint` yerel bilgisayarınızdaki (ana makine) `$HOME/sdh` dizinini konteyner içindeki `/home/tendermint` ev dizinine eşler.
* `-v $PWD:/apps` konteynerdeki iskeletlenmiş uygulamayı geçerli çalışma dizinindeki ana makineye kalıcı olarak yükler.
* `serve -p planet`, blockchain'in kaynak kodunu içeren `planet` dizinini kullanmayı belirtir.
* `-p 1317:1317`, konteyner içinde dinlenen 1317 numaralı portu ana makinedeki 1317 numaralı porta iletmek için API sunucu portunu (cosmos-sdk) ana makineyle eşler.
* `-p 26657:26657`, ana makinedeki 26657 (tendermint) RPC sunucu bağlantı noktasını Docker'daki 26657 bağlantı noktasıyla eşler.
* Blockchain başlatıldıktan sonra, Tendermint API'sini görmek için `http://localhost:26657` adresini açın.
* `-v` bayrağı, konteynerin uygulamanın kaynak koduna ana makineden erişmesini belirtir, böylece onu derleyebilir ve çalıştırabilir.

Docker konteynerinize hangi Ignite CLI sürümünün yükleneceğini ve çalıştırılacağını belirtebilirsiniz.

#### Son versiyon[​](broken-reference) <a href="#latest-version" id="latest-version"></a>

* Varsayılan olarak `ignite/cli`, `ignite/cli:latest` olarak çözümlenir.
* `latest` görüntü etiketi her zaman en son kararlı [Ignite CLI sürümüdür](https://github.com/ignite/cli/releases).

Örneğin, en son sürüm [v0.25.2](https://github.com/ignite/cli/releases/tag/v0.25.2) ise, en son etiketi `0.25.2` etiketine işaret eder.

#### Spesifik versiyon[​](broken-reference) <a href="#specific-version" id="specific-version"></a>

Ignite CLI'nin belirli bir sürümünü kullanmayı belirtebilirsiniz. Mevcut tüm etiketler Docker Hub'daki ignite/cli görüntüsündedir.

Örneğin:

* Sürüm `0.25.2`'yi kullanmak için `ignitehq/cli:0.25.2` (`v` öneki olmadan) kullanın.
* En son sürümü kullanmak için `ignitehq/cli` kullanın.
* `main` dalı (branch) kullanmak için `ignitehq/cli:main` kullanın, böylece gelecek sürümü deneyebilirsiniz.

En son görüntüyü almak için `docker pull`'u çalıştırın.

```
docker pull ignitehq/cli:main
```
