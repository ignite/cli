// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Ignite Cli Docs",
  tagline: "Ignite Cli Docs",
  url: "https://docs.ignite.com",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon-svg.svg",
  trailingSlash: false,

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "ignite",
  projectName: "ignite docs",

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  scripts: [
    {
      async: true,
      src: "https://www.googletagmanager.com/gtag/js?id=G-XL9GNV1KHW",
    },
  ],

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          versions: {
            current: {
              label: "nightly",
              path: "nightly",
              badge: true,
              banner: "unreleased", // put 'none' to remove
            },
          },
          sidebarPath: require.resolve("./sidebars.js"),
          routeBasePath: "/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      image: "img/og-image.jpg",
      announcementBar: {
        content:
          '<a target="_blank" rel="noopener noreferrer" href="https://ignite.com">← Back to Ignite</a>',
        isCloseable: false,
      },
      docs: {
        sidebar: {
          autoCollapseCategories: true,
        },
      },
      navbar: {
        hideOnScroll: true,
        logo: {
          alt: "Ignite Logo",
          src: "img/header-logo-docs.svg",
          srcDark: "img/header-logo-docs-dark.svg",
        },
        items: [
          {
            href: "https://github.com/ignite/cli",
            html: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="github-icon">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M12 0.300049C5.4 0.300049 0 5.70005 0 12.3001C0 17.6001 3.4 22.1001 8.2 23.7001C8.8 23.8001 9 23.4001 9 23.1001C9 22.8001 9 22.1001 9 21.1001C5.7 21.8001 5 19.5001 5 19.5001C4.5 18.1001 3.7 17.7001 3.7 17.7001C2.5 17.0001 3.7 17.0001 3.7 17.0001C4.9 17.1001 5.5 18.2001 5.5 18.2001C6.6 20.0001 8.3 19.5001 9 19.2001C9.1 18.4001 9.4 17.9001 9.8 17.6001C7.1 17.3001 4.3 16.3001 4.3 11.7001C4.3 10.4001 4.8 9.30005 5.5 8.50005C5.5 8.10005 5 6.90005 5.7 5.30005C5.7 5.30005 6.7 5.00005 9 6.50005C10 6.20005 11 6.10005 12 6.10005C13 6.10005 14 6.20005 15 6.50005C17.3 4.90005 18.3 5.30005 18.3 5.30005C19 7.00005 18.5 8.20005 18.4 8.50005C19.2 9.30005 19.6 10.4001 19.6 11.7001C19.6 16.3001 16.8 17.3001 14.1 17.6001C14.5 18.0001 14.9 18.7001 14.9 19.8001C14.9 21.4001 14.9 22.7001 14.9 23.1001C14.9 23.4001 15.1 23.8001 15.7 23.7001C20.5 22.1001 23.9 17.6001 23.9 12.3001C24 5.70005 18.6 0.300049 12 0.300049Z" fill="currentColor"/>
            </svg>
            `,
            position: "right",
          },
          {
            href: "https://ignite.com",
            className: "ignt-backlink",
            label: `Back to Ignite`,
            position: "right",
          },
          {
            type: "docsVersionDropdown",
            position: "left",
            dropdownActiveClassDisabled: true,
          },
        ],
      },
      footer: {
        links: [
          {
            items: [
              {
                html: `
                <a href="https://ignite.com">
                <svg width="83" height="25" viewBox="0 0 83 25" fill="none" xmlns="http://www.w3.org/2000/svg">
                <style>
                 path { fill: var(--ifm-font-color-base); }
                </style>
                <path d="M28.9089 4.71705C28.9089 5.61813 28.309 6.1741 27.4545 6.1741C26.6 6.1741 26 5.61813 26 4.71705C26 3.85433 26.6 3.27917 27.4545 3.27917C28.309 3.27917 28.9089 3.85433 28.9089 4.71705Z"/>
                <path d="M26.0359 19.0185H28.8773V7.57654H26.0359V19.0185Z" />
                <path d="M36.2042 24.3036C40.0947 24.3036 42.0618 22.0413 42.0618 18.9307V7.57575H39.2422V9.05494H39.0674C38.7177 8.48937 37.5811 7.29297 35.5485 7.29297C32.4885 7.29297 30.2373 9.64226 30.2373 13.2532C30.2373 16.8859 32.5323 19.1917 35.6359 19.1917C37.6904 19.1917 38.7614 17.9736 39.0892 17.4515H39.2641V18.9307C39.2641 20.8884 38.2587 21.9761 36.2042 21.9761C34.6523 21.9761 33.7344 21.3235 33.5158 20.3011H30.893C31.0897 22.4329 32.882 24.3036 36.2042 24.3036ZM36.2042 16.6466C34.3682 16.6466 33.1224 15.3415 33.1224 13.2532C33.1224 11.2302 34.3682 9.85979 36.1823 9.85979C37.9746 9.85979 39.286 11.1214 39.286 13.2532C39.286 15.1892 38.1057 16.6466 36.2042 16.6466Z" />
                <path d="M44.0167 19.0177H46.8581V12.4701C46.8581 11.0344 47.8416 9.9468 49.2186 9.9468C50.53 9.9468 51.4261 10.8822 51.4261 12.2743V19.0177H54.2675V11.5347C54.2675 9.09844 52.6938 7.29297 50.1147 7.29297C48.4755 7.29297 47.3389 8.11957 46.8799 8.96793H46.7269V7.57575H44.0167V19.0177Z" />
                <path d="M59.0895 4.71705C59.0895 5.61812 58.4895 6.1741 57.635 6.1741C56.7805 6.1741 56.1805 5.61812 56.1805 4.71705C56.1805 3.85433 56.7805 3.27917 57.635 3.27917C58.4895 3.27917 59.0895 3.85433 59.0895 4.71705Z"/>
                <path d="M56.2249 19.0185H59.0662V7.57654H56.2249V19.0185Z" />
                <path d="M64.8471 19.0189H67.0328V16.6696H65.4373C64.3444 16.6696 63.8636 16.0606 63.8636 15.1034V9.83928H67.0328V7.57699H63.8636V4.42285H61.0222V15.2122C61.0222 17.6703 62.4648 19.0189 64.8471 19.0189Z"/>
                <path d="M73.2954 19.3039C76.2679 19.3039 78.3225 17.7812 78.8907 15.4537H76.2461C75.8527 16.389 74.7817 17.0416 73.3829 17.0416C71.5469 17.0416 70.3448 15.867 70.2574 14.1267H79V13.2349C79 10.0155 76.9892 7.29639 73.2954 7.29639C69.9295 7.29639 67.5034 9.75444 67.5034 13.3654C67.5034 16.7806 69.9076 19.3039 73.2954 19.3039ZM70.3011 12.1255C70.4759 10.6898 71.6999 9.55867 73.2954 9.55867C74.9565 9.55867 76.1586 10.581 76.2898 12.1255H70.3011Z" />
                <path d="M12.0666 12.9609V18.7976C12.0667 18.8973 12.0407 18.9953 11.9911 19.0817C11.9415 19.1681 11.8701 19.2399 11.7841 19.2898L8.53588 21.087C8.49305 21.1121 8.44436 21.1254 8.39474 21.1255C8.34513 21.1256 8.29636 21.1126 8.25339 21.0878C8.21042 21.0629 8.17477 21.0271 8.15007 20.9839C8.12536 20.9408 8.11247 20.8919 8.11271 20.8421V15.1696C8.11272 15.0697 8.08649 14.9717 8.03666 14.8853C7.98684 14.7988 7.91517 14.7271 7.82889 14.6773L2.93106 11.8417C2.88804 11.8168 2.85232 11.781 2.82747 11.7379C2.80263 11.6948 2.78955 11.6459 2.78955 11.5962C2.78955 11.5464 2.80263 11.4975 2.82747 11.4544C2.85232 11.4113 2.88804 11.3755 2.93106 11.3507L6.17669 9.55226C6.26283 9.50254 6.36048 9.47638 6.45987 9.47638C6.55926 9.47638 6.65691 9.50254 6.74305 9.55226L11.7802 12.4699C11.8662 12.5197 11.9375 12.5913 11.9871 12.6775C12.0367 12.7637 12.0628 12.8614 12.0628 12.9609" />
                <path d="M19.9855 8.36439V17.6162C19.9853 17.7181 19.9583 17.8181 19.9074 17.9063C19.8564 17.9944 19.7832 18.0676 19.6952 18.1184L15.2551 20.6887C15.2111 20.7143 15.1613 20.7278 15.1105 20.7279C15.0597 20.7281 15.0098 20.7147 14.9658 20.6893C14.9217 20.6639 14.8852 20.6274 14.8597 20.5833C14.8343 20.5392 14.8209 20.4892 14.8209 20.4383V11.3543C14.8209 11.2524 14.7942 11.1522 14.7435 11.0638C14.6927 10.9754 14.6198 10.902 14.5318 10.8508L6.6848 6.30946C6.6408 6.28404 6.60425 6.24745 6.57883 6.20336C6.55342 6.15928 6.54004 6.10925 6.54004 6.05833C6.54004 6.00741 6.55342 5.95739 6.57883 5.9133C6.60425 5.86921 6.6408 5.83262 6.6848 5.8072L11.1314 3.23559C11.2196 3.18473 11.3195 3.15796 11.4211 3.15796C11.5228 3.15796 11.6227 3.18473 11.7108 3.23559L19.7004 7.86214C19.7885 7.91297 19.8616 7.98613 19.9126 8.07429C19.9636 8.16245 19.9905 8.2625 19.9907 8.36439" />
                </svg>
                </a>
              `,
              },
            ],
          },
          {
            title: "Products",
            items: [
              {
                label: "CLI",
                href: "https://ignite.com/cli",
              },
              {
                label: "Accelerator",
                href: "https://ignite.com/accelerator",
              },
            ],
          },
          {
            title: "Company",
            items: [
              {
                label: "About Ignite",
                href: "https://ignite.com/about",
              },
              {
                label: "Careers",
                href: "https://ignite.com/careers",
              },
              {
                label: "Blog",
                href: "https://ignite.com/blog",
              },
              {
                label: "Press",
                href: "https://ignite.com/press",
              },
            ],
          },
          {
            title: "Contact",
            items: [
              {
                label: "Business Inquiries",
                href: "mailto:business@ignite.com",
              },
            ],
          },
          {
            title: "Social",
            items: [
              {
                label: "Discord",
                href: "https://discord.com/invite/ignite",
              },
              {
                label: "Twitter",
                href: "https://twitter.com/ignite_com",
              },
              {
                label: "Linkedin",
                href: "https://www.linkedin.com/company/ignt/",
              },
              {
                label: "YouTube",
                href: "https://www.youtube.com/ignitehq",
              },
            ],
          },
        ],
        copyright: `<div>© Ignite ${new Date().getFullYear()}</div><div><a href="https://ignite.com/privacy">Privacy Policy</a><a href="https://ignite.com/terms-of-use">Terms of Use</a></div>`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ["protobuf", "go-module"], // https://prismjs.com/#supported-languages
        magicComments: [
          // Remember to extend the default highlight class name as well!
          {
            className: "theme-code-block-highlighted-line",
            line: "highlight-next-line",
            block: { start: "highlight-start", end: "highlight-end" },
          },
          {
            className: "code-block-removed-line",
            line: "remove-next-line",
            block: { start: "remove-start", end: "remove-end" },
          },
        ],
      },
      zoom: {
        selector: ".markdown :not(em) > img",
        config: {
          // options you can specify via https://github.com/francoischalifour/medium-zoom#usage
          background: {
            light: "rgb(255, 255, 255)",
            dark: "rgb(50, 50, 50)",
          },
        },
      },
      algolia: {
        appId: 'VVETP7QCVE',
        apiKey: '167213b8ce51cc7ff9a804df130657e5',
        indexName: 'ignite-cli',
        contextualSearch: true,

        // ↓ - To remove if `contextualSearch` versioning search works (to use if not)
        // exclusionPatterns: [
        //     'https://docs.ignite.com/v0.25.2/**',
        //     'https://docs.ignite.com/nightly/**',
        // ]
      },
    }),
  plugins: [
    [
      "@docusaurus/plugin-client-redirects",
      {
        createRedirects(existingPath) {
          if (existingPath.includes('/welcome')) {
            /*
            If the link received contains the path /guide, 
            this will change to /welcome.
            */ 
            return [
              existingPath.replace('/welcome', '/guide'),
            ];
          }
          return; // No redirect created if it doesn't contain /guide
        },
      },
    ],
    async function myPlugin(context, options) {
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
          postcssOptions.plugins.push(require("postcss-import"));
          postcssOptions.plugins.push(require("tailwindcss/nesting"));
          postcssOptions.plugins.push(require("tailwindcss"));
          postcssOptions.plugins.push(require("autoprefixer"));
          return postcssOptions;
        },
      };
    },
    require.resolve("docusaurus-plugin-image-zoom"),
  ],
};

module.exports = config;
