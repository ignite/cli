# TypeScript frontend

Ignite, blockchaininiz için istemci tarafı kodu oluşturmaya yönelik güçlü bir işlevsellik sunar. Bunu, blockchaininiz için özel olarak tasarlanmış tek tıklamayla istemci SDK oluşturma olarak düşünün.

TypeScript kod üretiminin nasıl kullanılacağı hakkında daha fazla bilgi edinmek için `ignite generate ts-client --help` bölümüne bakın.

`ignite scaffold chain` ile yeni bir blockchain oluşturun. Bunun yerine, varsa mevcut bir blockchain projesini kullanabilirsiniz.

```
ignite scaffold chain example
```

Test amacıyla `config.yml` dosyasına bir anımsatıcı ile yeni bir hesap ekleyin:

config.yml

```
accounts:
  - name: frank
    coins: ["1000token", "100000000stake"]
    mnemonic: play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint
```

Hem standart hem de özel Cosmos SDK modülleri için TypeScript istemcileri oluşturmak üzere bir komut çalıştırın:

```
ignite generate ts-client --clear-cache
```

Blockchain node'unuzu başlatmak için bir komut çalıştırın:

TypeScript istemcisi ile oluşturmaya başlamanın en iyi yolu [Vite ](https://vitejs.dev/)kullanmaktır. Vite, vanilla TS projelerinin yanı sıra React, Vue, Lit, Svelte ve Preact çerçeveleri için şablon kodu sağlar. [Vite Başlangıç kılavuzunda](https://vitejs.dev/guide) ek bilgi bulabilirsiniz.

Ayrıca istemcinin bağımlılıklarını çoklu doldurmanız ([polyfill](https://developer.mozilla.org/en-US/docs/Glossary/Polyfill)) gerekecektir. Aşağıda, gerekli çoklu doldurmalarla bir vanilya TS projesinin kurulumuna bir örnek verilmiştir:

```
npm create [email protected] my-frontend-app -- --template vanilla-ts
cd my-frontend-app
npm install --save-dev @esbuild-plugins/node-globals-polyfill @rollup/plugin-node-resolve
```

Daha sonra gerekli `vite.config.ts` dosyasını oluşturmalısınız.

my-frontend-app/vite.config.ts

```
import { nodeResolve } from "@rollup/plugin-node-resolve";
import { NodeGlobalsPolyfillPlugin } from "@esbuild-plugins/node-globals-polyfill";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [nodeResolve()],

  optimizeDeps: {
    esbuildOptions: {
      define: {
        global: "globalThis",
      },
      plugins: [
        NodeGlobalsPolyfillPlugin({
          buffer: true,
        }),
      ],
    },
  },
});
```

Daha sonra oluşturulan istemci kodunu bu proje içinde doğrudan veya istemciyi yayınlayarak ve diğer `npm`paketleri gibi yükleyerek kullanmaya hazırsınız.

Zincir başladıktan sonra Frank'in adresinin `cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7` olduğunu göreceksiniz. Frank'in hesabını bir sonraki bölümde verileri sorgulamak ve işlemleri yayınlamak için kullanacağız.

`ts-client`'ta oluşturulan kod, ihtiyaçlarınıza uyacak şekilde değiştirebileceğiniz yayınlamaya hazır bir `package.json` dosyasıyla birlikte gelir. `ts-client` için gerekli bağımlılıkları yükleyin:

İstemci, ihtiyacınız olan modülleri desteklemek ve örneklemek için bir istemci sınıfı yapılandırabileceğiniz modüler bir mimariye dayanmaktadır.

Varsayılan olarak, oluşturulan istemci, projenizde kullanılan tüm Cosmos SDK, özel ve 3. parti modülleri içeren bir istemci sınıfını dışa aktarır.

İstemciyi örneklemek için ortam bilgilerini (uç noktalar ve zincir öneki) sağlamanız gerekir. Sorgulama için ihtiyacınız olan tek şey bu:

my-frontend-app/src/main.ts

```
import { Client } from "../../ts-client";

const client = new Client(
  {
    apiURL: "http://localhost:1317",
    rpcURL: "http://localhost:26657",
    prefix: "cosmos",
  }
);
```

Yukarıdaki örnekte yerel bir dizindeki `ts-client` kullanılmıştır. Eğer `ts-client`'ınızı `npm`'de yayınladıysanız `../../ts-client` yerine bir paket adı yazın.

Sonuçta ortaya çıkan istemci örneği, her modül için bir `query` ve `tx` ad alanı ile modülün ilgili sorgulama ve işlem yöntemlerini tam tür ve otomatik tamamlama desteği ile içeren ad alanları içerir.

Bir adresin bakiyesini sorgulamak için:

```
const balances = await client.CosmosBankV1Beta1.query.queryAllBalances(
  'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7'
);
```

Bir anımsatıcıdan bir cüzdan oluşturarak (daha önce `config.yml` dosyasına eklenen Frank'in anımsatıcısını kullanıyoruz) ve bunu `Client()` işlevine isteğe bağlı bir argüman olarak aktararak istemciye imzalama özellikleri ekleyin. Cüzdan CosmJS OfflineSigner\` arayüzünü uygular.

my-frontend-app/src/main.ts

```
import { Client } from "../../ts-client";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic =
  "play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = new Client(
  {
    apiURL: "http://localhost:1317",
    rpcURL: "http://localhost:26657",
    prefix: "cosmos",
  },
  wallet
);
```

Bir işlemin yayınlanması:

my-frontend-app/src/main.ts

```
const tx_result = await client.CosmosBankV1Beta1.tx.sendMsgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
  fee: {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  memo: '',
})
```

Zincirinizde zaten tanımlanmış özel mesajlar varsa bunları kullanabilirsiniz. Eğer yoksa, örnek olarak Ignite'ın iskele kodunu kullanacağız. CRUD mesajları içeren bir gönderi oluşturun:

```
ignite scaffold list post title body
```

Zincirinize mesaj ekledikten sonra TypeScript istemcisini yeniden oluşturmanız gerekebilir:

```
ignite generate ts-client --clear-cache
```

Özel `MsgCreatePost` içeren bir işlem yayınlayın:

my-frontend-app/src/main.ts

```
import { Client } from "../../ts-client";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic =
  "play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = new Client(
  {
    apiURL: "http://localhost:1317",
    rpcURL: "http://localhost:26657",
    prefix: "cosmos",
  },
  wallet
);
const tx_result = await client.ExampleExample.tx.sendMsgCreatePost({
  value: {
    title: 'foo',
    body: 'bar',
    creator: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
  },
  fee: {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  memo: '',
})
```

İsterseniz, genel istemci sınıfını içe aktararak ve ihtiyacınız olan modüllerle genişleterek yalnızca ilgilendiğiniz modülleri kullanarak daha hafif bir istemci oluşturabilirsiniz:

my-frontend-app/src/main.ts

```
import { IgniteClient } from '../../ts-client/client'
import { Module as CosmosBankV1Beta1 } from '../../ts-client/cosmos.bank.v1beta1'
import { Module as CosmosStakingV1Beta1 } from '../../ts-client/cosmos.staking.v1beta1'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'

const mnemonic =
  'play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint'
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)
const Client = IgniteClient.plugin([CosmosBankV1Beta1, CosmosStakingV1Beta1])

const client = new Client(
  {
    apiURL: 'http://localhost:1317',
    rpcURL: 'http://localhost:26657',
    prefix: 'cosmos',
  },
  wallet,
)
```

Ayrıca TX mesajlarını ayrı ayrı oluşturabilir ve bunları aşağıdaki gibi bir global imzalama istemcisi kullanarak tek bir TX içinde gönderebilirsiniz:

my-frontend-app/src/main.ts

```
const msg1 = await client.CosmosBankV1Beta1.tx.msgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
})

const msg2 = await client.CosmosBankV1Beta1.tx.msgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
})

const tx_result = await client.signAndBroadcast(
  [msg1, msg2],
  {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  '',
)
```

Son olarak, daha fazla kullanım kolaylığı için, yukarıda bahsedilen modüler istemcinin yanı sıra, oluşturulan her modül ayrı bir txClient ve queryClient'ı açığa çıkararak sadeleştirilmiş bir şekilde kendi başına kullanılabilir.

my-frontend-app/src/main.ts

```
import { txClient } from '../../ts-client/cosmos.bank.v1beta1'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'

const mnemonic =
  'play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint'
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic)

const client = txClient({
  signer: wallet,
  prefix: 'cosmos',
  addr: 'http://localhost:26657',
})

const tx_result = await client.sendMsgSend({
  value: {
    amount: [
      {
        amount: '200',
        denom: 'token',
      },
    ],
    fromAddress: 'cosmos13xkhcx2dquhqdml0k37sr7yndquwteuvt2cml7',
    toAddress: 'cosmos15uw6qpxqs6zqh0zp3ty2ac29cvnnzd3qwjntnc',
  },
  fee: {
    amount: [{ amount: '0', denom: 'stake' }],
    gas: '200000',
  },
  memo: '',
})
```

Normalde Keplr, `OfflineSigner` arayüzünü uygulayan bir cüzdan nesnesi sağlar, böylece istemci örneklemesindeki `wallet` argümanını `window.keplr.getOfflineSigner(chainId)` ile değiştirebilirsiniz. Bununla birlikte, Keplr zinciriniz hakkında zincir kimliği, denomlar, ücretler vb. gibi bilgilere ihtiyaç duyar. [`experimentalSuggestChain()`](https://docs.keplr.app/api/suggest-chain.html), Keplr'ın bu bilgileri Keplr uzantısına iletmek için sağladığı bir yöntemdir.

Oluşturulan istemci, zincir bilgilerini otomatik olarak keşfeden ve sizin için ayarlayan bir `useKeplr()`yöntemi sunarak bunu kolaylaştırır. Böylece, istemciyi bir cüzdan olmadan örnekleyebilir ve ardından Keplr aracılığıyla işlem yapmayı etkinleştirmek için `useKeplr()` yöntemini çağırabilirsiniz:

my-frontend-app/src/main.ts

```
import { Client } from '../../ts-client';

const client = new Client({ 
        apiURL: "http://localhost:1317",
        rpcURL: "http://localhost:26657",
        prefix: "cosmos"
    }
);
await client.useKeplr();
```

`useKeplr()` isteğe bağlı olarak, otomatik olarak keşfedilen değerleri geçersiz kılmanıza olanak tanıyan `experimentalSuggestChain()` işlevinin `ChainInfo` türü bağımsız değişkeniyle aynı anahtarlardan bir veya daha fazlasını içeren bir nesne bağımsız değişkenini kabul eder.

Örneğin, varsayılan zincir adı ve token hassasiyeti (zincir üzerinde kaydedilmez) `<chainId> Network`ve `0` olarak ayarlanırken, denom için ticker büyük harfle denom adına ayarlanır. Bunları geçersiz kılmak istiyorsanız, şöyle bir şey yapabilirsiniz:

my-frontend-app/src/main.ts

```
import { Client } from '../../ts-client';

const client = new Client({ 
        apiURL: "http://localhost:1317",
        rpcURL: "http://localhost:26657",
        prefix: "cosmos"
    }
);
await client.useKeplr({
  chainName: 'My Great Chain',
  stakeCurrency: {
    coinDenom: 'TOKEN',
    coinMinimalDenom: 'utoken',
    coinDecimals: '6',
  },
})
```

İstemci ayrıca, halihazırda başlatılmış bir istemcide cüzdanı farklı bir cüzdanla değiştirmenize de olanak tanır:

```
import { Client } from '../../ts-client';
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

const mnemonic =
  'play butter frown city voyage pupil rabbit wheat thrive mind skate turkey helmet thrive door either differ gate exhibit impose city swallow goat faint'
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic);

const client = new Client({ 
        apiURL: "http://localhost:1317",
        rpcURL: "http://localhost:26657",
        prefix: "cosmos"
    }
);
await client.useKeplr();

// broadcast transactions using the Keplr wallet

client.useSigner(wallet);

// broadcast transactions using the CosmJS wallet
```
