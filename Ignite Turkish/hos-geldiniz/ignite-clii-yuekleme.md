# Ignite CLI'ı Yükleme

[Ignite CLI](https://github.com/ignite/cli)'yı web tabanlı bir Gitpod IDE'de çalıştırabilir veya Ignite CLI'yı yerel bilgisayarınıza yükleyebilirsiniz.

### Ön Koşullar

Ignite CLI'yı yüklemeden ve kullanmadan önce ön koşulları karşıladığınızdan emin olun.&#x20;

#### İşletim sistemleri

Ignite CLI aşağıdaki işletim sistemleri için desteklenir:

* GNU/Linux
* macOS
* Linux için Windows Alt Sistemi (WSL)

### Go

Ignite CLI, Go programlama dilinde yazılmıştır. Ignite CLI'yı yerel bir sistemde kullanmak için:

* [Go](https://golang.org/doc/install)'yu yükleyin (**sürüm 1.19** veya üstü)
* Go ortam değişkenlerinin sisteminizde [doğru şekilde](https://golang.org/doc/gopath\_code#GOPATH) ayarlandığından emin olun

### Ignite CLI sürümünüzü doğrulayın

Yüklediğiniz Ignite CLI sürümünü doğrulamak için aşağıdaki komutu çalıştırın:

```
ignite version
```

### Ignite CLI'ı Yükleme <a href="#installing-ignite-cli" id="installing-ignite-cli"></a>

Ignite ikili dosyasının en son sürümünü yüklemek için aşağıdaki komutu kullanın:

```
curl https://get.ignite.com/cli! | bash
```

Bu komut, kurulum betiğini indirmek için curl'ü çağırır ve kurulumu gerçekleştirmek için çıktıyı bash'e aktarır. ignite binary'si `/usr/local/bin` dosyasına yüklenir.

Daha fazla bilgi edinmek veya yükleme işlemini özelleştirmek için GitHub'daki [yükleyici belgelerine](https://github.com/ignite/installer) bakın.

### Yazma izni

Ignite CLI kurulumu `/usr/local/bin/` dizinine yazma izni gerektirir. Eğer `/usr/local/bin/` dizinine yazma izniniz olmadığı için kurulum başarısız olursa, aşağıdaki komutu çalıştırın:

```
curl https://get.ignite.com/cli | bash
```

Ardından ignite çalıştırılabilir dosyasını `/usr/local/bin/` konumuna taşımak için bu komutu çalıştırın:

```
sudo mv ignite /usr/local/bin/
```

Bazı makinelerde bir izin hatası oluşur:

```
mv: rename ./ignite to /usr/local/bin/ignite: Permission denied
============
Error: mv failed
```

Bu durumda, sudo'yu curl'den önce ve bash'ten önce kullanın:

```
sudo curl https://get.ignite.com/cli | sudo bash
```

### Ignite CLI kurulumunuzun yükseltilmesi

Ignite CLI'ın yeni bir sürümünü yüklemeden önce, mevcut tüm Ignite CLI yüklemelerini kaldırın.

Mevcut Ignite CLI kurulumunu kaldırmak için:

* Terminal pencerenizde, ignite chain serve ile başlattığınız zinciri durdurmak için Ctrl+C tuşlarına basın.
* Ignite CLI binary dosyasını rm $(which ignite) ile kaldırın. Kullanıcı izinlerinize bağlı olarak, komutu sudo ile veya sudo olmadan çalıştırın.
* Tüm ignite kurulumları sisteminizden kaldırılana kadar bu adımı tekrarlayın.

Mevcut tüm Ignite CLI kurulumları kaldırıldıktan sonra, [Ignite CLI Kurulumu](https://docs.ignite.com/welcome/install#installing-ignite-cli) talimatlarını izleyin.

Sürüm özellikleri ve değişiklikleri hakkında ayrıntılı bilgi için depodaki [changelog.md](https://github.com/ignite/cli/blob/main/changelog.md) dosyasına bakın.

### Kaynaktan derleme

Kaynak kodu denemek için kaynaktan derleme yapabilirsiniz:

```
git clone https://github.com/ignite/cli --depth=1
cd cli && make install
```

### Özet

* Ön koşulları doğrulayın.
* Yerel bir geliştirme ortamı kurmak için Ignite CLI'yi bilgisayarınıza yerel olarak yükleyin.
* Ignite CLI'yı cURL kullanarak ikili dosyayı getirerek veya kaynaktan oluşturarak yükleyin.
* Varsayılan olarak en son sürüm yüklenir. Önceden derlenmiş ignite binary'sinin önceki sürümlerini yükleyebilirsiniz.
* Yeni bir sürüm yüklemeden önce zinciri durdurun ve mevcut sürümleri kaldırın.
