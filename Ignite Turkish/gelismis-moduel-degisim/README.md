# Gelişmiş Modül: Değişim

Interchain Exchange, blok zincirleri arasında alım ve satım emirleri oluşturmak için kullanılan bir modüldür.

Bu eğitimde, emir çiftleri, alım emirleri ve satış emirleri oluşturabilen bir Cosmos SDK modülünün nasıl oluşturulacağını öğreneceksiniz. Blok zincirleri arasında emir defterleri ve alım satım emirleri oluşturursunuz, bu da bir blok zincirinden diğerine token takası yapmanızı sağlar.

**Not**: Bu eğitimdeki kod özellikle bu eğitim için yazılmıştır ve yalnızca eğitim amaçlıdır. Bu eğitim kodunun üretimde kullanılması amaçlanmamıştır.

Sonucu görmek istiyorsanız, [değişim reposundaki](https://github.com/tendermint/interchange) örnek uygulamaya bakın.

#### **Nasıl yapılacağını öğrenecekleriniz**:

* Ignite CLI ile bir blok zinciri oluşturma
* Cosmos SDK IBC modülü oluşturma
* Bir modül ile alış ve satış emirlerini barındıran bir emir defteri oluşturma
* IBC paketlerini bir blok zincirinden diğerine gönderme
* IBC paketlerinin zaman aşımları ve onayları ile ilgilenme

### Değişim Değişim Modülü Nasıl Çalışır?

İki veya daha fazla blok zinciri ile çalışan bir borsa oluşturmak için, `dex` adlı bir Cosmos SDK modülü oluşturmak üzere bu eğitimdeki adımları izleyin.

Yeni `dex` modülü, bir çift token için bir borsa emir defteri açmanıza olanak tanır: bir blok zincirinden bir token ve başka bir blok zincirindeki bir token. Blok zincirlerinin dex modülüne sahip olması gerekmektedir.

Token, basit bir emir defterinde limit emirleri ile satın alınabilir veya satılabilir. Bu eğitimde, likidite havuzu veya otomatik piyasa yapıcı (AMM) kavramı yoktur.

Piyasa tek yönlüdür:

* Kaynak zincirde satılan token geri alınamaz çünkü
* Hedef zincirden satın alınan token aynı çift kullanılarak geri satılamaz.

Kaynak zincirdeki bir token satılırsa, yalnızca emir defterinde yeni bir çift oluşturularak geri alınabilir. Bu iş akışı, hedef blok zincirinde bir `voucher` tokeni oluşturan Blok Zincirler Arası İletişim protokolünün (IBC) doğasından kaynaklanmaktadır. Yerel bir blok zinciri tokenı ile başka bir blok zincirinde basılan bir `voucher` tokenı arasında fark vardır. Yerel tokenı geri almak için ikinci bir sipariş defteri çifti oluşturmanız gerekir.

Bir sonraki bölümde, blok zincirleri arası değişimin tasarımıyla ilgili ayrıntıları öğreneceksiniz.
