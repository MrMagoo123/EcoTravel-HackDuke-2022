// hook the open and change url
const proxied = window.XMLHttpRequest.prototype.open;
console.log("pls save us")

// shape security pt2
window.XMLHttpRequest.prototype.open = function () {
  console.log("XMLHttpRequest hook")
  arguments[1] = arguments[1].replace("https://bl2751.wixsite.com", document.location.protocol + "//" + document.location.host)

  return proxied.apply(this, [].slice.call(arguments));
}

if (window.viewerModel && window.viewerModel.site) {
  window.viewerModel.site.externalBaseUrl = document.location.protocol + "//" + document.location.host + "/";
} else {
  window.viewerModel = {
    site: {}
  }
}

Object.defineProperty(window.viewerModel.site, "externalBaseUrl", {
  value: document.location.protocol + "//" + document.location.host + "/",
  writable: false,
  configurable: false,
});

Object.freeze(window.viewerModel.site);