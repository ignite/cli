# DeFi Loan

Merkezi olmayan finans (DeFi), blok zinciri ekosisteminin hızla büyüyen ve finansal araçlar ve hizmetler hakkındaki düşüncelerimizi dönüştüren bir sektörüdür. DeFi, internet bağlantısı ve dijital cüzdanı olan herkesin erişebileceği borç verme, borç alma, spot ticaret, marj ticareti ve flaş krediler dahil olmak üzere çok çeşitli yenilikçi finansal ürün ve hizmetler sunmaktadır.

DeFi'nin en önemli avantajlarından biri, son kullanıcıların karmaşık işe alım süreçlerine veya pasaport veya geçmiş kontrolleri gibi kişisel belgelerin sunulmasına gerek kalmadan finansal araçlara ve hizmetlere hızlı ve kolay bir şekilde erişmelerine olanak sağlamasıdır. Bu da DeFi'yi yavaş, maliyetli ve zahmetli olabilen geleneksel bankacılık sistemlerine karşı cazip bir alternatif haline getiriyor.

Bu eğitimde, kullanıcıların birbirlerine dijital varlık ödünç vermelerini ve birbirlerinden ödünç almalarını sağlayan bir DeFi platformunun nasıl oluşturulacağını öğreneceksiniz. Oluşturacağınız platform, tüm işlemlerin merkezi olmayan ve değişmez bir kaydını sağlayan bir blok zinciri tarafından desteklenecektir. Bu, platformun şeffaf, güvenli ve dolandırıcılığa karşı dirençli olmasını sağlar.

Kredi, bir tarafın, borç alanın, para veya dijital token gibi belirli bir miktar varlık aldığı ve kredi tutarını artı bir ücreti önceden belirlenmiş bir son tarihe kadar borç verene geri ödemeyi kabul ettiği finansal bir işlemdir. Krediyi güvence altına almak için, borçlu teminat sağlar ve bu teminat, borçlunun krediyi kararlaştırıldığı şekilde geri ödememesi durumunda borç veren tarafından ele geçirilebilir.

Bir kredinin hüküm ve koşullarını tanımlayan çeşitli özellikleri vardır.

`ID`, krediyi blok zincirinde tanımlamak için kullanılan benzersiz bir tanımlayıcıdır.

`amount`, borçluya ödünç verilen varlıkların miktarıdır.

`fee`, borçlunun kredi için borç verene ödemesi gereken maliyettir.

`collateral`, borçlunun kredi için teminat olarak borç verene sağladığı varlık veya varlıklardır.

`deadline`, borçlunun krediyi geri ödemesi gereken tarihtir. Borçlu krediyi son ödeme tarihine kadar geri ödeyemezse, borç veren krediyi tasfiye etmeyi ve teminata el koymayı seçebilir.

Bir kredinin `state`'i, kredinin mevcut durumunu tanımlar ve `requested`, `approved`, `paid`, `cancelled` veya `liquidated` gibi çeşitli değerler alabilir. Bir kredi, borçlu kredi için ilk kez bir talep gönderdiğinde `requested` durumdadır. Borç veren talebi onaylarsa, kredi `approved` durumuna geçer. Borçlu krediyi geri ödediğinde, kredi `paid` durumuna geçer. Borçlu kred-iyi onaylanmadan önce iptal ederse, kredi `cancelled` duruma geçer. Borçlu krediyi son ödeme tarihine kadar geri ödeyemezse, borç veren krediyi tasfiye etmeyi ve teminata el koymayı seçebilir. Bu durumda, kredi `liquidated` duruma geçer.

Bir kredi işleminde iki taraf vardır: borç alan ve borç veren. Borçlu, krediyi talep eden ve kredi tutarını artı bir ücreti önceden belirlenmiş bir son tarihe kadar kredi verene geri ödemeyi kabul eden taraftır. Kredi veren, kredi talebini onaylayan ve borçluya kredi tutarını sağlayan taraftır.

Bir borçlu olarak, kredi platformunda çeşitli eylemler gerçekleştirebilmeniz gerekir. Bu eylemler şunları içerebilir:

* kredi talebinde bulunmak,
* Bir krediyi iptal etmek,
* Bir kredinin geri ödenmesi.

Bir kredi talebinde bulunmak, miktar, ücret, teminat ve geri ödeme için son tarih dahil olmak üzere kredinin hüküm ve koşullarını belirlemenize olanak tanır. Bir krediyi iptal ederseniz, kredi onaylanmadan veya finanse edilmeden önce kredi talebinizi geri çekebilirsiniz. Bir kredinin geri ödenmesi, kredi tutarını artı ücreti kredi koşullarına uygun olarak kredi verene geri ödemenizi sağlar.

Bir borç veren olarak, platformda iki eylem gerçekleştirebilmeniz gerekir:

* bir krediyi onaylamak
* bir krediyi tasfiye etmek.

Bir krediyi onaylamak, kredinin hüküm ve koşullarını kabul etmenize ve kredi tutarını borçluya göndermenize olanak tanır. Bir kredinin tasfiye edilmesi, krediyi son ödeme tarihine kadar geri ödeyememeniz durumunda borç verenin teminata el koymasına olanak tanır.

Bu eylemleri gerçekleştirerek, borç verenler ve borç alanlar birbirleriyle etkileşime girebilir ve platformdaki dijital varlıkların ödünç verilmesini ve ödünç alınmasını kolaylaştırabilir. Platform, kullanıcıların varlıklarını yönetmelerine ve finansal hedeflerine güvenli ve şeffaf bir şekilde ulaşmalarına olanak tanıyan finansal araçlara ve hizmetlere erişmelerini sağlar.
