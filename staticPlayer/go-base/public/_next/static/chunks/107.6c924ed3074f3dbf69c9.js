(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[107],{4107:function(e,n,t){"use strict";t.r(n);var r=t(5893),o=t(1163),i=t(4865),s=t.n(i),a=t(7294);s().configure({showSpinner:!1,minimum:.1,trickleSpeed:150});n.default=function(){var e=(0,a.useRef)(void 0),n=(0,a.useState)(!1),t=n[0],i=n[1],u=function(){t||(i(!0),e.current=setTimeout((function(){s().start()}),250))},c=function(){i(!1),clearTimeout(e.current),s().done()};return(0,a.useEffect)((function(){return o.default.events.on("routeChangeStart",u),o.default.events.on("routeChangeComplete",c),o.default.events.on("routeChangeError",c),function(){o.default.events.off("routeChangeStart",u),o.default.events.off("routeChangeComplete",c),o.default.events.off("routeChangeError",c)}})),t?(0,r.jsx)("div",{className:"progressParOverlay"}):null}},4865:function(e,n,t){var r,o;void 0===(o="function"===typeof(r=function(){var e={version:"0.2.0"},n=e.settings={minimum:.08,easing:"ease",positionUsing:"",speed:200,trickle:!0,trickleRate:.02,trickleSpeed:800,showSpinner:!0,barSelector:'[role="bar"]',spinnerSelector:'[role="spinner"]',parent:"body",template:'<div class="bar" role="bar"><div class="peg"></div></div><div class="spinner" role="spinner"><div class="spinner-icon"></div></div>'};function t(e,n,t){return e<n?n:e>t?t:e}function r(e){return 100*(-1+e)}function o(e,t,o){var i;return(i="translate3d"===n.positionUsing?{transform:"translate3d("+r(e)+"%,0,0)"}:"translate"===n.positionUsing?{transform:"translate("+r(e)+"%,0)"}:{"margin-left":r(e)+"%"}).transition="all "+t+"ms "+o,i}e.configure=function(e){var t,r;for(t in e)void 0!==(r=e[t])&&e.hasOwnProperty(t)&&(n[t]=r);return this},e.status=null,e.set=function(r){var a=e.isStarted();r=t(r,n.minimum,1),e.status=1===r?null:r;var u=e.render(!a),c=u.querySelector(n.barSelector),f=n.speed,l=n.easing;return u.offsetWidth,i((function(t){""===n.positionUsing&&(n.positionUsing=e.getPositioningCSS()),s(c,o(r,f,l)),1===r?(s(u,{transition:"none",opacity:1}),u.offsetWidth,setTimeout((function(){s(u,{transition:"all "+f+"ms linear",opacity:0}),setTimeout((function(){e.remove(),t()}),f)}),f)):setTimeout(t,f)})),this},e.isStarted=function(){return"number"===typeof e.status},e.start=function(){e.status||e.set(0);var t=function(){setTimeout((function(){e.status&&(e.trickle(),t())}),n.trickleSpeed)};return n.trickle&&t(),this},e.done=function(n){return n||e.status?e.inc(.3+.5*Math.random()).set(1):this},e.inc=function(n){var r=e.status;return r?("number"!==typeof n&&(n=(1-r)*t(Math.random()*r,.1,.95)),r=t(r+n,0,.994),e.set(r)):e.start()},e.trickle=function(){return e.inc(Math.random()*n.trickleRate)},function(){var n=0,t=0;e.promise=function(r){return r&&"resolved"!==r.state()?(0===t&&e.start(),n++,t++,r.always((function(){0===--t?(n=0,e.done()):e.set((n-t)/n)})),this):this}}(),e.render=function(t){if(e.isRendered())return document.getElementById("nprogress");u(document.documentElement,"nprogress-busy");var o=document.createElement("div");o.id="nprogress",o.innerHTML=n.template;var i,a=o.querySelector(n.barSelector),c=t?"-100":r(e.status||0),f=document.querySelector(n.parent);return s(a,{transition:"all 0 linear",transform:"translate3d("+c+"%,0,0)"}),n.showSpinner||(i=o.querySelector(n.spinnerSelector))&&l(i),f!=document.body&&u(f,"nprogress-custom-parent"),f.appendChild(o),o},e.remove=function(){c(document.documentElement,"nprogress-busy"),c(document.querySelector(n.parent),"nprogress-custom-parent");var e=document.getElementById("nprogress");e&&l(e)},e.isRendered=function(){return!!document.getElementById("nprogress")},e.getPositioningCSS=function(){var e=document.body.style,n="WebkitTransform"in e?"Webkit":"MozTransform"in e?"Moz":"msTransform"in e?"ms":"OTransform"in e?"O":"";return n+"Perspective"in e?"translate3d":n+"Transform"in e?"translate":"margin"};var i=function(){var e=[];function n(){var t=e.shift();t&&t(n)}return function(t){e.push(t),1==e.length&&n()}}(),s=function(){var e=["Webkit","O","Moz","ms"],n={};function t(e){return e.replace(/^-ms-/,"ms-").replace(/-([\da-z])/gi,(function(e,n){return n.toUpperCase()}))}function r(n){var t=document.body.style;if(n in t)return n;for(var r,o=e.length,i=n.charAt(0).toUpperCase()+n.slice(1);o--;)if((r=e[o]+i)in t)return r;return n}function o(e){return e=t(e),n[e]||(n[e]=r(e))}function i(e,n,t){n=o(n),e.style[n]=t}return function(e,n){var t,r,o=arguments;if(2==o.length)for(t in n)void 0!==(r=n[t])&&n.hasOwnProperty(t)&&i(e,t,r);else i(e,o[1],o[2])}}();function a(e,n){return("string"==typeof e?e:f(e)).indexOf(" "+n+" ")>=0}function u(e,n){var t=f(e),r=t+n;a(t,n)||(e.className=r.substring(1))}function c(e,n){var t,r=f(e);a(e,n)&&(t=r.replace(" "+n+" "," "),e.className=t.substring(1,t.length-1))}function f(e){return(" "+(e.className||"")+" ").replace(/\s+/gi," ")}function l(e){e&&e.parentNode&&e.parentNode.removeChild(e)}return e})?r.call(n,t,n,e):r)||(e.exports=o)}}]);