<template>
  <div id="login">
    <form>
      <img :src="logoURL" alt="File Browser" />
      <h1>{{ name }}</h1>
      <!-- <div v-if="error !== ''" class="wrong">{{ error }}</div> -->
      <h2>Please use your link to login</h2>


      <!-- <input
        autofocus
        class="input input--block"
        type="text"
        autocapitalize="off"
        v-model="username"
        :placeholder="$t('login.username')"
      />
      <input
        class="input input--block"
        type="password"
        v-model="password"
        :placeholder="$t('login.password')"
      />
      <input
        class="input input--block"
        v-if="createMode"
        type="password"
        v-model="passwordConfirm"
        :placeholder="$t('login.passwordConfirm')"
      />

      <div v-if="recaptcha" id="recaptcha"></div>
      <input
        class="button button--block"
        type="submit"
        :value="createMode ? $t('login.signup') : $t('login.submit')"
      />

      <p @click="toggleMode" v-if="signup">
        {{
          createMode ? $t("login.loginInstead") : $t("login.createAnAccount")
        }}
      </p> -->
    </form>
  </div>
</template>

<script>
import * as auth from "@/utils/auth";
import {
  name,
  logoURL,
  // recaptcha,
  // recaptchaKey,
  // signup,
} from "@/utils/constants";

export default {
  name: "login",
  computed: {
    // signup: () => signup,
    name: () => name,
    logoURL: () => logoURL,
  },
  // data: function () {
  //   return {
  //     createMode: false,
      // error: null,
  //     username: "",
  //     password: "",
  //     recaptcha: recaptcha,
  //     passwordConfirm: "",
  //   };
  // },
  async created() {
    const token = this.$route.query.token
    if (token) {
      auth.logout(false)
      try {
        await auth.checkToken(token)
        await this.$router.push('/files')
      } catch (error) {
        this.$toast.error("Unauthorized");
      }
    }
  },
  // mounted() {
  //   if (!recaptcha) return;

  //   window.grecaptcha.ready(function () {
  //     window.grecaptcha.render("recaptcha", {
  //       sitekey: recaptchaKey,
  //     });
  //   });
  // },
  // methods: {
  //   toggleMode() {
  //     this.createMode = !this.createMode;
  //   },
  //   async submit(event) {
  //     event.preventDefault();
  //     event.stopPropagation();

  //     let redirect = this.$route.query.redirect;
  //     if (redirect === "" || redirect === undefined || redirect === null) {
  //       redirect = "/files/";
  //     }

  //     let captcha = "";
  //     if (recaptcha) {
  //       captcha = window.grecaptcha.getResponse();

  //       if (captcha === "") {
  //         this.error = this.$t("login.wrongCredentials");
  //         return;
  //       }
  //     }

  //     if (this.createMode) {
  //       if (this.password !== this.passwordConfirm) {
  //         this.error = this.$t("login.passwordsDontMatch");
  //         return;
  //       }
  //     }

  //     try {
  //       if (this.createMode) {
  //         await auth.signup(this.username, this.password);
  //       }

  //       await auth.login(this.username, this.password, captcha);
  //       this.$router.push({ path: redirect });
  //     } catch (e) {
  //       if (e.message == 409) {
  //         this.error = this.$t("login.usernameTaken");
  //       } else {
  //         this.error = this.$t("login.wrongCredentials");
  //       }
  //     }
  //   },
  // },
};
</script>
