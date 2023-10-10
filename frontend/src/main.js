import "whatwg-fetch";
import cssVars from "css-vars-ponyfill";
import { sync } from "vuex-router-sync";
import store from "@/store";
import router from "@/router";
import i18n from "@/i18n";
import Vue from "@/utils/vue";
// import { recaptcha } from "@/utils/constants";
import App from "@/App";

cssVars();

sync(store, router);

async function start() {
  // if (recaptcha) {
  //   await new Promise((resolve) => {
  //     const check = () => {
  //       if (typeof window.grecaptcha === "undefined") {
  //         setTimeout(check, 100);
  //       } else {
  //         resolve();
  //       }
  //     };

  //     check();
  //   });
  // }

  new Vue({
    el: "#app",
    store,
    router,
    i18n,
    template: "<App/>",
    components: { App },
  });
}

start();
