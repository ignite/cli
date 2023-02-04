# Dokümantasyonun iyileştirilmesi

Depomuzu ziyaret ettiğiniz ve katkıda bulunmayı düşündüğünüz için teşekkür ederiz. Harika öğreticiler ve belgeler oluşturmamıza ve sürdürmemize yardımcı olmak için gösterdiğiniz ilgiye minnettarız.

Sorunuzun daha önce sorulmuş ve yanıtlanmış olup olmadığını görmek için mevcut [Ignite CLI sorunlarını](https://github.com/ignite/cli/issues) inceleyin.

* Geri bildirim sağlamak, bir sorun bildirmek ve nasıl daha iyi hale getirebileceğimizi anlamamıza yardımcı olacak cömert ayrıntılar sağlamak için.
* Bir düzeltme sağlamak için doğrudan katkıda bulunun. Üye veya bakımcı değilseniz, depoyu çatallayın ve ardından çatallanmış deponuzdan `main` branch'e bir çekme isteği (PR) gönderin.
* Taslak bir çekme isteği oluşturarak başlayın. Çalışmanız yeni başlıyor veya tamamlanmamış olsa bile taslak PR'nizi erkenden oluşturun. Taslak PR'niz topluluğa bir şey üzerinde çalıştığınızı gösterir ve geliştirme sürecinin başlarında konuşmalar için bir alan sağlar. `Draft` PR'ler için birleştirme engellenmiştir, bu nedenle deneme yapmak ve yorum davet etmek için güvenli bir yer sağlarlar.

### Teknik içerik PR'larının gözden geçirilmesi <a href="#reviewing-technical-content-prs" id="reviewing-technical-content-prs"></a>

En iyi içerik katkılarından bazıları PR inceleme döngüleri sırasında gelir. Tıpkı kod incelemelerinde yaptığınız gibi teknik içerik PR incelemeleri için de en iyi uygulamaları takip edin.

* Satır içi öneriler için [GitHub önerme özelliğini](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/commenting-on-a-pull-request) kullanın.
* PR sahibi, önerdiğiniz değişiklikleri teker teker veya toplu olarak (tercih edilir) birleştirebilir.
* 20'den fazla satır içi öneriyle sonuçlanan daha ayrıntılı bir kapsamlı inceleme sağlıyorsanız, devam edin ve dalı kontrol edin ve değişiklikleri kendiniz yapın.

Belgelere ve öğreticilere katkıları memnuniyetle karşılıyoruz.

Teknik içeriğimiz [Google geliştirici belgeleri stil kılavuzunu](https://developers.google.com/style) takip etmektedir. Başlamanıza yardımcı olacak önemli noktalar:

* [Öne Çıkanlar](https://developers.google.com/style/highlights)
* [Kelime Listesi](https://developers.google.com/style/word-list)
* [Stil ve ton](https://developers.google.com/style/tone)
* [Global bir kitle için yazmak](https://developers.google.com/style/translation)
* [Çapraz referanslar](https://developers.google.com/style/cross-references)
* [Şimdiki zaman](https://developers.google.com/style/tense)

Google yönergeleri burada listelenenden daha fazla materyal içerir ve önerilen içerik değişiklikleri hakkında kolay karar vermeyi sağlayan bir kılavuz olarak kullanılır.

Diğer faydalı kaynaklar:

* [Google Teknik Yazı Yazma Kursları](https://developers.google.com/tech-writing)
* [GitHub Kılavuzları Markdown'da Uzmanlaşma](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)

Teknik içerik, bilgi tabanı makalelerini ve etkileşimli eğitimleri içerir.

* Ignite CLI Developer Tutorials içeriği `docs/guide` klasöründedir.
* Bilgi Tabanı içeriği `docs/kb` klasöründedir.
* Yükseltme bilgileri `docs/migration` klasöründe bulunmaktadır.

Note: The CLI docs are auto-generated and do not support doc updates.

Diğer içeriklerin konumları ve klasörleri değişiklik gösterebilir. İlgilendiğiniz içerik için kendini tanımlayan klasörleri keşfedin. Bazı makaleler ve eğitimler tek bir Markdown dosyasında bulunurken, diğer eğitimler için alt klasörler mevcut olabilir.

Her zaman olduğu gibi, üzerinde çalışılmakta olan içerikler başka konumlarda ve depolarda olabilir.

Ignite ürün ekibi geliştiricileri, Ignite CLI'yi oluşturmaya ve geliştirici deneyimini iyileştirmeye odaklanmıştır. Ignite Ekosistem Geliştirme ekibi, teknik içeriğin ve eğitimlerin sahibidir ve geliştirici katılımını yönetir.

[Ignite CLI'nin arkasındaki kişilerle ve katkıda bulunanlarımızla](https://github.com/ignite/cli/graphs/contributors) tanışın.

Güncellenen sayfalar yayınlanmadan önce değişikliklerinizin üretimde nasıl görüneceğini görmek için bir önizleme kullanın.

* Bir PR taslak modundayken, Markdown'daki önizleme özelliğini kullanmaya güvenebilirsiniz.
* PR **Taslak**'tan **İnceleme için Hazır**'a geçtikten sonra, CI durum kontrolleri bir dağıtım önizlemesi oluşturur. Siz çalışmaya ve aynı dala yeni değişiklikler işlemeye devam ettikçe bu önizleme güncel kalır. Bir GitHub eylemleri URL'sindeki `Docs Deploy Preview / build_and_deploy (pull_request)` önizlemesi o PR için benzersizdir.
