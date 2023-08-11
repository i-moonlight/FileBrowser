<template>
  <nav :class="{ active }">
    <template v-if="isLogged">
      <button
        class="action"
        @click="toRoot"
        :aria-label="$t('sidebar.myFiles')"
        :title="$t('sidebar.myFiles')"
      >
        <i class="material-icons">folder</i>
        <span>{{ $t("sidebar.myFiles") }}</span>
      </button>

      <div v-if="user.perm.create">
        <button
          @click="$store.commit('showHover', 'newDir')"
          class="action"
          :aria-label="$t('sidebar.newFolder')"
          :title="$t('sidebar.newFolder')"
        >
          <i class="material-icons">create_new_folder</i>
          <span>{{ $t("sidebar.newFolder") }}</span>
        </button>

        <button
          @click="$store.commit('showHover', 'newFile')"
          class="action"
          :aria-label="$t('sidebar.newFile')"
          :title="$t('sidebar.newFile')"
        >
          <i class="material-icons">note_add</i>
          <span>{{ $t("sidebar.newFile") }}</span>
        </button>
      </div>

      <div>
        <button
          @click="logout"
          class="action"
          id="logout"
          :aria-label="$t('sidebar.logout')"
          :title="$t('sidebar.logout')"
        >
          <i class="material-icons">exit_to_app</i>
          <span>{{ $t("sidebar.logout") }}</span>
        </button>
      </div>
    </template>
    <template v-else>
      <router-link
        class="action"
        to="/login"
        :aria-label="$t('sidebar.login')"
        :title="$t('sidebar.login')"
      >
        <i class="material-icons">exit_to_app</i>
        <span>{{ $t("sidebar.login") }}</span>
      </router-link>
    </template>
  </nav>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import * as auth from "@/utils/auth";
import {
  version,
  disableExternal,
  disableUsedPercentage,
} from "@/utils/constants";

export default {
  name: "sidebar",
  computed: {
    ...mapState(["user"]),
    ...mapGetters(["isLogged"]),
    active() {
      return this.$store.state.show === "sidebar";
    },
    version: () => version,
    disableExternal: () => disableExternal,
    disableUsedPercentage: () => disableUsedPercentage,
  },
  methods: {
    toRoot() {
      this.$router.push({ path: "/files/" }, () => {});
      this.$store.commit("closeHovers");
    },
    toSettings() {
      this.$router.push({ path: "/settings" }, () => {});
      this.$store.commit("closeHovers");
    },
    help() {
      this.$store.commit("showHover", "help");
    },
    logout: auth.logout,
  },
};
</script>
