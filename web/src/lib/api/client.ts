import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  signup(body: api.SignupBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signup", "POST", api.Signup, z.any(), body, options)
  }
  
  signin(body: api.SigninBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signin", "POST", api.Signin, z.any(), body, options)
  }
  
  changePassword(body: api.ChangePasswordBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/password", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  startDownload(options?: ExtraOptions) {
    return this.request("/api/v1/system/download", "GET", z.undefined(), z.any(), undefined, options)
  }
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  importMalList(body: api.ImportMalListBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/import/mal", "POST", z.undefined(), z.any(), body, options)
  }
  
  getAnimes(options?: ExtraOptions) {
    return this.request("/api/v1/animes", "GET", api.GetAnimes, z.any(), undefined, options)
  }
  
  getAnimeById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/animes/${id}`, "GET", api.GetAnimeById, z.any(), undefined, options)
  }
  
  setAnimeUserData(id: string, body: api.SetAnimeUserData, options?: ExtraOptions) {
    return this.request(`/api/v1/animes/${id}/user`, "POST", z.undefined(), z.any(), body, options)
  }
}
