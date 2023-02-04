# Ignite Zinciri

_Ignite, Cosmos SDK tabanlı blockchain'lerin başlatılmasına yardımcı olacak bir blockchain'dir._

Cosmos SDK ve Ignite CLI kullanarak, geliştiriciler hızlı bir şekilde merkezi olmayan, kullanım için ekonomik ve ölçeklenebilir bir kripto uygulaması oluşturabilirler. Cosmos SDK çerçevesi, geliştiricilerin daha geniş [Cosmos ekosisteminin](https://v1.cosmos.network/ecosystem/apps) bir parçası haline gelen egemen uygulamaya özgü blockchain'ler oluşturmasına olanak tanır. Cosmos SDK ile oluşturulan blockchainler, zinciri güvence altına almak için doğrulayıcılar gerektiren bir Proof-of-Stake (PoS) konsensüs protokolü kullanır.

Ignite CLI gibi araçlar bir Cosmos SDK blockchain'inin geliştirilmesini basitleştirse de, yeni bir zincir başlatmak oldukça karmaşık bir süreçtir. Kendi egemen blockchain'inizi geliştirmenin ve başlatmanın en büyük zorluklarından biri, altta yatan fikir birliğinin güvenliğini sağlamaktır. Cosmos SDK zincirleri PoS mutabakatına dayandığından, her bir blockchain başlatılmadan önce ilk coin tahsislerine ve validatörlere ihtiyaç duyar, bu da geliştiricilere zincirlerinin tokenomiklerini belirlemek veya sağlam bir validatör setini koordine etmek gibi önemli zorluklar sunar.

İlk jeton tahsisleri ve validatörler, ağdaki tüm ilk düğümler arasında paylaşılan JSON formatlı bir genesis dosyasında tanımlanır. Bu genesis dosyası uygulamanın başlangıç durumunu tanımlar. PoS'a dayanan güvenli zincirler, coinlerin ilk tahsisinin iyi bir şekilde dağıtılmasını gerektirir, böylece hiçbir validatör tüm tokenların 1/3'ünden fazlasına sahip olamaz ve orantısız miktarda oylama gücü elde edemez.

Temel konsensüs güvenliğini sağlamanın yanı sıra, yeni bir blockchain başlatmanın bir diğer zor görevi de genesis dosyası için çeşitli doğrulayıcıların ilgisini çekmektir. Gelecek vaat eden birçok proje, kaynak veya deneyim eksikliği nedeniyle zincirlerini güvence altına almak için yeterli sayıda güvenilir validatörün dikkatini çekmeyi başaramamaktadır.

Bu nedenle Ignite Zinciri, geliştiricilerin bir blockchain başlatmanın karmaşıklıklarında gezinmelerine ve yeni bir zincirin oluşumunu koordine etmelerine yardımcı olarak Cosmos SDK blockchain'lerinin başlatılmasını kolaylaştırmak için tasarlandı. Blockchain'in merkezi olmayan doğasını kullanan Ignite'ın koordinasyon özellikleri, blockchain kurucularının validatörler ve yatırımcılarla bağlantı kurmasına yardımcı olarak projelerinin piyasaya sürülme süresini ve başarı şansını hızlandırır.

Ignite Chain ile etkileşim komutları Ignite CLI'ya entegre edilmiştir ve zincirlerin buradan başlatılmasına izin verir. Ignite Chain ile entegrasyon, CLI'nin geliştiriciyi, blockchain'in geliştirilmesi ve denenmesinden ana ağının başlatılmasına kadar bir Cosmos projesini gerçekleştirmenin tüm yaşam döngüsünde desteklemesine olanak tanır.

Ignite Chain, Cosmos SDK tabanlı zincirlerin lansmanını basitleştiren, koordinasyon, hazırlık ve lansman aşamalarında hayati kaynaklar ve destek sağlayan güvenli bir platformdur. Ignite, blockchain projelerinin, validatör koordinasyonu ve token ihracından bağış toplama ve topluluk oluşturmaya kadar zincirlerini başlatmanın karmaşıklıklarının üstesinden gelmek için ihtiyaç duydukları araçları sağlar.

Ignite, üç aşamalı genel bir lansman süreci ile yeni zincirlerin launch'ını kolaylaştırır:

* Koordinasyon
* Hazırlık
* Launch

Her aşamada sürtünmeyi azaltmak için Ignite, validatör koordinasyonu için değişmez ve evrensel bir veritabanı sağlar.

Gelecekte, Ignite ayrıca şunları da sunacaktır:

* Token ihracı: Ignite, gelecekteki bir mainnet'in hisse tahsisini temsil eden tokenlerin (kupon olarak adlandırılır) çıkarılmasına izin verir
* Kupon satışı için bir bağış toplama platformu
* Başlatılan bir testnet ağında validatör faaliyetlerini ödüllendirmek için izinsiz bir çerçeve

Cosmos ekosisteminde bir zincir başlatmak için, validatörler yeni blok zinciri ağını oluşturmak üzere birbirlerine bağlanan nodeları başlatmalıdır. Bir node, genesis dosyası adı verilen bir dosyadan başlatılmalıdır. Yeni zincirin başlatılabilmesi için genesis dosyasının tüm doğrulayıcı düğümlerde aynı olması gerekir.

JSON biçimli genesis dosyası, coin tahsisleri, validatörlerin listesi, blokları aktif olarak imzalayan maksimum validatör sayısı gibi zincir için çeşitli parametreler ve belirli başlatma zamanı dahil olmak üzere zincirin başlangıç durumu hakkında bilgi içerir. Her validatör aynı genesis dosyasına sahip olduğundan, genesis zamanına ulaşıldığında blok zinciri ağı otomatik olarak başlar.

#### Gerçeğin koordinasyon kaynağı olarak Ignite[​](broken-reference) <a href="#ignite-as-a-coordination-source-of-truth" id="ignite-as-a-coordination-source-of-truth"></a>

Ignite Chain, yeni zincirlerin bir validatör setini koordine etmesi ve validatörlerin bir zincir launch'ı için genesis oluşturması için bir doğruluk kaynağı görevi görür. Blockchain, nihai genesis dosyasını doğrudan kendi defterinde saklamaz, bunun yerine genesis dosyasının deterministik bir şekilde oluşturulmasına olanak tanıyan bilgileri saklar.

Ignite'ta depolanan ve belirli bir zincir başlatma için genesis dosyasının deterministik olarak oluşturulmasını destekleyen bilgiler _launch bilgileri_ olarak adlandırılır. Ignite üzerinde yeni bir zincir oluşturulurken koordinatör ilk başlatma bilgisini sağlar. Daha sonra, zincir içi koordinasyon yoluyla, bu başlatma bilgileri mesajlar gönderilerek blockchain ile etkileşime girilerek güncellenir. Zincir başlatılmaya hazır olduğunda, başlatma bilgilerini kullanan bir genesis oluşturma algoritması çağrılarak genesis dosyası oluşturulur.

**GenesisGenerate(LaunchInformation) => genesis.json**

Genesis oluşturma algoritması resmi ve resmi olarak belirtilmiştir. Genesis oluşturma algoritmasının resmi uygulaması Ignite CLI kullanılarak Go'da geliştirilmiştir. Bununla birlikte, herhangi bir proje, algoritmanın spesifikasyonuna uyduğu sürece algoritmanın kendi uygulamasını geliştirmekte özgürdür.

Genesis oluşturma algoritması zincir içi protokolün bir parçası değildir. Yeni bir zinciri başarıyla başlatmak için, tüm validatörlerin başlatma bilgilerini kullanarak kendi genesislerini oluşturmak için algoritmayı kullanmaları gerekir. Algoritma, Ignite zincirinde depolanan başlatma bilgilerinden deterministik olarak oluşumu üretir.

Başlatma bilgilerinin herhangi bir unsuru sansürlenirse, örneğin bir hesap bakiyesi kaldırılırsa, başlatılan zincir itibarı olumsuz etkilenir ve validatörlerin çoğunluğunun kullanmama konusunda hemfikir olduğu anlamına gelir:

* Kurcalamaya dayanıklı launch bilgileri
* Resmi genesis oluşturma algoritması

Genesis oluşturma dışında, genesis oluşturma algoritması spesifikasyonu ağ yapılandırmanızı nasıl kuracağınız konusunda rehberlik eder. Örneğin, başlatma bilgileri blockchain ağının kalıcı eşlerinin adreslerini içerebilir.

Başlatma bilgileri üç farklı şekilde oluşturulabilir ya da güncellenebilir:

1. Zincir oluşturma sırasında tanımlanır ancak oluşturulduktan sonra koordinatör tarafından güncellenebilir
2. Koordinasyon yoluyla belirlenir
3. Koordinasyonla ilgili olmayan belirli zincir içi mantık aracılığıyla belirlenir

#### 1 - Zincir oluşturma sırasında belirlenen başlatma bilgileri. <a href="#1---launch-information-determined-during-chain-creation" id="1---launch-information-determined-during-chain-creation"></a>

* `GenesisChainID`: Ağ için tanımlayıcı
* `SourceURL`: Blockchain node binary'sini oluşturmak için kullanılan kaynak kodun git deposunun URL'si
* `SourceHash`: Kaynak kodun sürümünü tanımlayan özel hash
* `InitialGenesis`: Genesis oluşturma algoritmasını çalıştırmadan önce zincir başlatma için başlangıç genesis'ini belirten çok formatlı bir yapı

#### 2 - Koordinasyon yoluyla belirlenen launch bilgileri. <a href="#2---launch-information-determined-through-coordination" id="2---launch-information-determined-through-coordination"></a>

* `GenesisAccounts`: Zincir için, ilişkili bakiyeleri olan adreslerden oluşan bir oluşum hesapları listesi
* `VestingAccounts`: Hakediş seçenekleri olan genesis hesaplarının listesi
* `GenesisValidators`: Zincir başlatıldığında ilk validatörlerin bir listesi
* `ParamChanges`: Genesis durumundaki modül param değişikliklerinin bir listesi

#### 3 - Zincir içi mantık yoluyla belirlenen launch bilgileri. <a href="#3---launch-information-determined-through-on-chain-logic" id="3---launch-information-determined-through-on-chain-logic"></a>

* `GenesisTime`: Ağ başlangıcı için LaunchTime olarak da adlandırılan zaman damgası

#### İlk genesis[​](broken-reference) <a href="#initial-genesis" id="initial-genesis"></a>

Başlatma bilgileri başlangıç genesis yapısını içerir. Bu yapı, genesis oluşturma algoritmasını çalıştırmadan ve genesis dosyasını sonlandırmadan önce ilk genesisin oluşturulması için bilgi sağlar.

Başlangıç genesis yapısı şunlar olabilir:

* `DefaultGenesis`: varsayılan genesis dosyası chain binary init komutu tarafından oluşturulur
* `GenesisURL`: bir zincir başlatma için ilk oluşum, bir URL'den alınan ve daha sonra gerekli algoritma ile değiştirilen mevcut bir oluşum dosyasıdır - bu ilk oluşum türü, ilk oluşum durumu kapsamlı olduğunda, token dağıtımı için çok sayıda hesap içerdiğinde, bir airdrop için kayıtlar içerdiğinde kullanılmalıdır
* `GenesisConfig`: Bir zincir başlatma için ilk oluşum, oluşum hesaplarını ve modül parametrelerini içeren bir Ignite CLI yapılandırmasından oluşturulur - bu ilk oluşum türü, koordinatörün ilk oluşum için kapsamlı bir durumu olmadığında ancak bazı modül parametrelerinin özelleştirilmesi gerektiğinde kullanılmalıdır. Örneğin, staking token için staking bond denom

Koordinasyon süreci zincir oluşturulduktan hemen sonra başlar ve koordinatör zincirin başlatılmasını tetiklediğinde sona erer.

Başlatma bilgileri koordinasyon süreci sırasında güncellenir.

Koordinasyon süreci sırasında herhangi bir varlık ağa istek gönderebilir. İstek, içeriği başlatma bilgilerindeki güncellemeleri belirten bir nesnedir.

Zincir koordinatörü talepleri onaylar veya reddeder:

* Bir talep onaylanırsa, içerik launch bilgilerine uygulanır
* Talep reddedilirse, başlatma bilgilerinde herhangi bir değişiklik yapılmaz

İstek oluşturucu ayrıca isteği doğrudan reddedebilir veya iptal edebilir.

Her zincir, tüm istekleri içeren bir istek havuzu içerir. Her isteğin bir durumu vardır:

* _PENDING_: Koordinatörün onayının beklenmesi
* _APPROVED_: Koordinatör tarafından onaylanmış, içeriği lansman bilgilerine uygulanmıştır
* _REJECTED_: Koordinatör veya talep oluşturucu tarafından reddedildi

Approving or rejecting a request is irreversible. The only possible status transitions are:

* _PENDING_ 'den _APPROVED 'a_
* _PENDING_ 'den _REJECTED 'e_

Bir talebin başlatma bilgileri üzerindeki etkisini geri almak için, bir kullanıcı nihai karşıt talebi göndermelidir (örnek: AddAccount → RemoveAccount).

Koordinatör talepler için tek onaylayıcı olduğundan, koordinatör tarafından oluşturulan her talep derhal ONAYLANDI olarak ayarlanır ve içeriği başlatma bilgilerine uygulanır.

Ignite zincirine altı tür istek gönderilebilir:

* `AddGenesisAccount`
* `AddVestingAccount`
* `AddGenesisValidator`
* `RemoveAccount`
* `RemoveValidator`
* `ChangeParam`

**`AddGenesisAccount`**genesis zinciri için coin bakiyesi olan yeni bir hesap talep eder. Bu istek içeriği iki alandan oluşur:

* Hesap adresi, launch bilgilerinde benzersiz olmalıdır
* Hesap bakiyesi

Lansman bilgilerinde zaten aynı adrese sahip bir genesis hesabı veya bir vesting hesabı belirtilmişse, talep otomatik olarak uygulanamaz.

**`AddVestingAccount`** genesis zinciri için bir coin bakiyesi ve hak ediş seçenekleri ile yeni bir hesap talep eder. Bu istek içeriği iki alandan oluşur:

* Hesabın adresi
* Hesabın hak ediş seçenekleri

Şu anda desteklenen hak ediş seçeneği, hesabın toplam bakiyesinin belirtildiği ve hesabın toplam bakiyesinin belirli sayıda jetonunun yalnızca bir bitiş zamanına ulaşıldıktan sonra hak edildiği gecikmeli hak ediştir.

Lansman bilgilerinde bir genesis hesabı veya aynı adrese sahip bir hak ediş hesabı zaten belirtilmişse, talep otomatik olarak uygulanamaz.

**`AddGenesisValidator`** zincir için yeni bir genesis validatörü talep eder. Cosmos SDK blok zincirindeki bir genesis validatörü, ağ başladığında bağlı bir validatör olmak için genesis başlatma sırasında bakiyesinin bir kısmını kendi kendine devreden genesis'te mevcut bir bakiyeye sahip bir hesabı temsil eder. Çoğu durumda, doğrulayıcı, zincirin ilk oluşumunda bakiyesi olan bir hesaba zaten sahip değilse, doğrulayıcı olmayı talep etmeden önce `AddGenesisAccount`ile bir hesap talep etmelidir.

Genesis başlatma sırasında kendi kendine delegasyon, [genutils adlı bir Cosmos SDK modülü](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/genutil) ile gerçekleştirilir. Genesis'te _genutils_ modülü, ağ başlatılmadan önce yürütülen işlemleri temsil eden gentx adlı nesneleri içerir. Ağ başladığında bir validatör olmak için, gelecekteki bir validatörün kendi hesabından self-delegation işlemini içeren bir gentx sağlaması gerekir.

Talep içeriği beş alandan oluşur:

* Validatör öz delegasyonu için gentx
* Validatörün adresi
* Validatör node'un konsensüs açık anahtarı
* Kendi kendine delegasyon
* Validatör node için eş bilgisi

Aynı adrese sahip bir validatör başlatma bilgilerinde zaten mevcutsa, istek otomatik olarak uygulanamaz.

**`RemoveAccount`** bir genesis veya vesting hesabının lansman bilgilerinden kaldırılmasını talep eder. İstek içeriği kaldırılacak hesabın adresini içerir. Launch bilgilerinde belirtilen adrese sahip bir genesis veya vesting hesabı yoksa istek otomatik olarak uygulanamaz.

**`RemoveValidator`** bir genesis validatörünün launch bilgilerinden kaldırılmasını talep eder. İstek içeriği kaldırılacak validatörün adresini içerir. Launch bilgilerinde belirtilen adrese sahip bir validatör hesabı yoksa istek otomatik olarak uygulanamaz.

**`ChangeParam`** genesis'teki bir modül parametresinin değiştirilmesini talep eder. Cosmos SDK blok zincirindeki modüller, blockchain mantığını yapılandıracak parametrelere sahip olabilir. Parametreler, blockchain ağı canlı olduğunda yönetişim yoluyla değiştirilebilir. Başlatma işlemi sırasında, zincirin ilk parametreleri genesis'te ayarlanır.

Bu istek içeriği üç alandan oluşur:

* Modülün adı
* Parametrenin adı
* Jenerik veri olarak temsil edilen parametrenin değeri

#### Geçerlilik talebi <a href="#request-validity" id="request-validity"></a>

Bir talep uygulanırken bazı kontroller zincir üzerinde doğrulanır. Örneğin, bir genesis hesabı iki kez eklenemez. Ancak, diğer bazı geçerlilik özellikleri zincir üzerinde kontrol edilemez. Örneğin, bir gentx blockchain'de genel bir bayt dizisi ile temsil edildiğinden, gentx'in doğru bir şekilde imzalandığını veya zincir üzerinde depolanan sağlanan konsensüs ortak anahtarının gentx'teki konsensüs ortak anahtarına karşılık geldiğini doğrulamak için zincir üzerinde bir kontrol mümkün değildir. Bu gentx doğrulaması, taleplerin geçerli bir formata sahip olmasını ve zincirin başlamasına izin vermesini sağlamak için blockchain ile etkileşime giren istemcinin sorumluluğundadır. Genesis oluşturma algoritmasında bazı geçerlilik kontrolleri belirtilmiştir.

Ignite aracılığıyla bir zincirin genel başlatma(launch) süreci üç aşamadan oluşur:

* Koordinasyon aşaması&#x20;
* Hazırlık aşaması&#x20;
* Launch aşaması

Koordinatör Ignite üzerinde zinciri oluşturduktan ve ilk başlatma bilgilerini sağladıktan sonra başlatma süreci, kullanıcıların zincir oluşumu için talep gönderebilecekleri koordinasyon aşamasına girer. Koordinatör zincirin başlatılmaya hazır olduğunu düşündükten sonra zincirin başlatılmasını tetikler. Bu işlem sırasında koordinatör, zincir için başlatma zamanını veya genesis zamanını sağlar.

Fırlatma tetiklendikten sonra ve fırlatma süresine ulaşılmadan önce, zincir fırlatma süreci hazırlık aşamasına girer. Hazırlık aşamasında artık istek gönderilemez ve zincirin başlatma bilgileri son haline getirilir. Validatörler, zincirin nihai genesis'ini almak için genesis oluşturma algoritmasını çalıştırır ve nodelerini hazırlar. Kalan süre, validatörlerin nodelarını hazırlamaları için yeterli zamanı sağlamalıdır. Bu başlatma süresi koordinatör tarafından belirlenir, ancak kalan süre için belirli bir aralık uygulanır.

Başlatma zamanına ulaşıldığında zincir ağı başlatılır ve zincir başlatma süreci başlatma aşamasına girer. Bu noktada zincir canlı olduğu için koordinatörün başka bir işlem yapmasına gerek yoktur. Ancak bazı durumlarda zincir başlatılamamış olabilir. Örneğin, oluşumdaki her validatör kendi node'unu başlatmazsa zincir başlamaz.

Koordinatör zincirleme fırlatmayı geri döndürme yeteneğine sahiptir. Zincir başlatmanın geri alınması, başlatma sürecini, başlatma başarısızlığıyla ilgili sorunun ele alınmasına izin vermek için isteklerin yeniden gönderilebileceği koordinasyon aşamasına geri döndürür. Başlatmanın geri alınmasının yalnızca Ignite üzerinde etkisi vardır. Yeni zincir etkin bir şekilde başlatılırsa, Ignite'ta başlatmanın geri alınmasının zincirin canlılığı üzerinde hiçbir etkisi yoktur. Zincirin başlatılmasının geri döndürülmesi yalnızca başlatma süresi artı geri döndürme gecikmesi adı verilen bir gecikmeden sonra koordinatör tarafından gerçekleştirilebilir.

Determinizmi sağlamak için, zincirin launch bilgisine bağlı olarak oluşum oluşturma kuralları titizlikle belirlenmelidir.&#x20;

Oluşum üretimi için genel adımlar şunlardır:

* Kaynaktan blockchain node binary'sini oluşturmak&#x20;
* İlk genesisin oluşturulması&#x20;
* Zincir kimliğinin ayarlanması&#x20;
* Oluşum zamanının ayarlanması&#x20;
* Genesis hesapları ekleme&#x20;
* Hakediş seçenekli genesis hesapları ekleme&#x20;
* Genesis doğrulayıcıları için gentxs ekleme&#x20;
* Param değişikliklerinden modül parametrelerini değiştirme
