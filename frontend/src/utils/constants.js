const name = window.FileBrowser.Name || "File Browser";
const disableExternal = window.FileBrowser.DisableExternal;
const disableUsedPercentage = window.FileBrowser.DisableUsedPercentage;
const baseURL = window.FileBrowser.BaseURL;
const staticURL = window.FileBrowser.StaticURL;
const recaptcha = window.FileBrowser.ReCaptcha;
const recaptchaKey = window.FileBrowser.ReCaptchaKey;
const version = window.FileBrowser.Version;
const logoURL = `${staticURL}/img/logo.png`;
const noAuth = window.FileBrowser.NoAuth;
const loginPage = window.FileBrowser.LoginPage;
const theme = window.FileBrowser.Theme;
const enableThumbs = window.FileBrowser.EnableThumbs;
const resizePreview = window.FileBrowser.ResizePreview;
const enableExec = window.FileBrowser.EnableExec;
const origin = window.location.origin;

export {
  name,
  disableExternal,
  disableUsedPercentage,
  baseURL,
  logoURL,
  recaptcha,
  recaptchaKey,
  version,
  noAuth,
  loginPage,
  theme,
  enableThumbs,
  resizePreview,
  enableExec,
  origin,
};
