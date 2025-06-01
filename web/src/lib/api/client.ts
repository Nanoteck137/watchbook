import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, createUrl, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  url: ClientUrls;

  constructor(baseUrl: string) {
    super(baseUrl);
    this.url = new ClientUrls(baseUrl);
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
  
  importMalAnime(body: api.ImportMalAnimeBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/import/mal/anime", "POST", api.ImportMalAnime, z.any(), body, options)
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
  
  getUserAnimeList(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/animes/user/list/${id}`, "GET", api.GetAnimes, z.any(), undefined, options)
  }
  
}

export class ClientUrls {
  baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }
  
  signup() {
    return createUrl(this.baseUrl, "/api/v1/auth/signup")
  }
  
  signin() {
    return createUrl(this.baseUrl, "/api/v1/auth/signin")
  }
  
  changePassword() {
    return createUrl(this.baseUrl, "/api/v1/auth/password")
  }
  
  getMe() {
    return createUrl(this.baseUrl, "/api/v1/auth/me")
  }
  
  getSystemInfo() {
    return createUrl(this.baseUrl, "/api/v1/system/info")
  }
  
  startDownload() {
    return createUrl(this.baseUrl, "/api/v1/system/download")
  }
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/system/sse")
  }
  
  updateUserSettings() {
    return createUrl(this.baseUrl, "/api/v1/user/settings")
  }
  
  importMalList() {
    return createUrl(this.baseUrl, "/api/v1/user/import/mal")
  }
  
  importMalAnime() {
    return createUrl(this.baseUrl, "/api/v1/user/import/mal/anime")
  }
  
  getAnimes() {
    return createUrl(this.baseUrl, "/api/v1/animes")
  }
  
  getAnimeById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/animes/${id}`)
  }
  
  setAnimeUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/animes/${id}/user`)
  }
  
  getUserAnimeList(id: string) {
    return createUrl(this.baseUrl, `/api/v1/animes/user/list/${id}`)
  }
  
  getAnimeImage(id: string, image: string) {
    return createUrl(this.baseUrl, `/files/animes/${id}/${image}`)
  }
}
