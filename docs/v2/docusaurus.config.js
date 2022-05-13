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
  favicon: "img/favicon.ico",
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

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
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
      navbar: {
        logo: {
          alt: "My Site Logo",
          src: "img/logo.svg",
        },
        items: [
          {
            type: "doc",
            docId: "index",
            position: "left",
            label: "Docs",
          },
          {
            href: "https://github.com/ignite-hq/cli",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            items: [
              {
                html: `
                <a href="https://ignite.com">
                  <img src="img/logo.svg" alt="ignite logo" width="114" height="51" />
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
              {
                label: "Ventures",
                href: "https://ignite.com/ventures",
              },
              {
                label: "Emeris",
                href: "https://emeris.com",
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
                label: "Media Inquiries",
                href: "mailto:media@tendermint.com",
              },
              {
                label: "Business Inquiries",
                href: "mailto:business@tendermint.com",
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
            ],
          },
        ],
        copyright: `© Ignite ${new Date().getFullYear()}`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
  plugins: [
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
  ],
};

module.exports = config;
