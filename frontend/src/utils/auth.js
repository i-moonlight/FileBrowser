import store from "@/store";
import router from "@/router";
import { Base64 } from "js-base64";
import { baseURL } from "@/utils/constants";

export function parseToken(token, sessionId) {
  const parts = token.split(".");

  if (parts.length !== 3) {
    throw new Error("token malformed");
  }

  const data = JSON.parse(Base64.decode(parts[1]));

  document.cookie = `auth=${token}; path=/`;

  localStorage.setItem("jwt", token);
  store.commit("setJWT", token);
  store.commit("setSessionId", sessionId);
  store.commit("setUser", data.user);
}

export async function checkToken(jwt, sessionId) {
  const res = await fetch(`${baseURL}/api/check-token`, {
    method: 'POST',
    headers: { 'X-Auth': jwt, 'X-Session-Id': sessionId },
  });

  if (res.status === 200) {
    parseToken(jwt, sessionId);
  } else {
    throw new Error(res);
  }
}

export async function mount(jwt, sessionId) {
  const res = await fetch(`${baseURL}/api/mount`, {
    method: "POST",
    headers: { "X-Auth": jwt, 'X-Session-Id': sessionId },
  });

  if (res.status === 200) {
    return
  } else {
    throw new Error(res);
  }
}

export function logout(isRedirect = true) {
  document.cookie = "auth=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/";

  store.commit("setJWT", "");
  store.commit("setUser", null);
  localStorage.setItem("jwt", null);
  if (isRedirect) {
    router.push({ path: "/login" });
  }
}
