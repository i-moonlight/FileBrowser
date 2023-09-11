<template>
  <main id="login">
    <div v-if="loading">
      <h2 class="message delayed">
        <div class="spinner">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span>{{ $t("files.loading") }}</span>
      </h2>
    </div>
    
    <form v-else>
      <img :src="logoUrl" alt="File Browser" />
      <h1>{{ name }}</h1>
      <h2>{{ $t("Please use your link to login") }}</h2>
    </form>
  </main>
</template>

<script>
import * as auth from "@/utils/auth";
import { mapState, mapMutations } from "vuex";
import { name, logoURL } from "@/utils/constants";

export default {
  name: "login",
  computed: {
    ...mapState(["loading"]),
    name: () => name,
    logoUrl: () => {
      return logoURL
    }
  },
  async created() {
    const token = this.$route.query.token
    let redirect = this.$route.query.redirect;

    if (redirect === "" || redirect === undefined || redirect === null) {
      redirect = "/files/";
    }

    if (token) {
      this.setLoading(true)
      auth.logout(false)

      const sessionId = crypto.randomUUID()
      sessionStorage.setItem('token', token)
      sessionStorage.setItem('sessionId', sessionId)

      try {
        await auth.checkToken(token, sessionId)
        await auth.mount(token, sessionId)
        await this.$router.push(redirect)
        this.$toast.success("Welcome")
      } catch (error) {
        this.$toast.error(this.$t("Unauthorized"))
      } finally {
        this.setLoading(false)
      }
    }
  },
  methods: {
    ...mapMutations(["setLoading"]),
  }
};
</script>
