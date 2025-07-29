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
  
  syncLibrary(options?: ExtraOptions) {
    return this.request("/api/v1/system/library", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  cleanupLibrary(options?: ExtraOptions) {
    return this.request("/api/v1/system/library/cleanup", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  getAllApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "GET", api.GetAllApiTokens, z.any(), undefined, options)
  }
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getMedia(options?: ExtraOptions) {
    return this.request("/api/v1/media", "GET", api.GetMedia, z.any(), undefined, options)
  }
  
  getMediaById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}`, "GET", api.GetMediaById, z.any(), undefined, options)
  }
  
  getMediaParts(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts`, "GET", api.GetMediaParts, z.any(), undefined, options)
  }
  
  setMediaUserData(id: string, body: api.SetMediaUserData, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/user`, "POST", z.undefined(), z.any(), body, options)
  }
  
  getCollections(options?: ExtraOptions) {
    return this.request("/api/v1/collections", "GET", api.GetCollections, z.any(), undefined, options)
  }
  
  getCollectionById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "GET", api.GetCollectionById, z.any(), undefined, options)
  }
  
  getCollectionItems(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items`, "GET", api.GetCollectionItems, z.any(), undefined, options)
  }
  
  createCollection(body: api.CreateCollectionBody, options?: ExtraOptions) {
    return this.request("/api/v1/collections", "POST", api.CreateCollection, z.any(), body, options)
  }
  
  editCollection(id: string, body: api.EditCollectionBody, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  providerMyAnimeListGetAnime(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/provider/myanimelist/anime/${id}`, "GET", api.ProviderMyAnimeListAnime, z.any(), undefined, options)
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
  
  syncLibrary() {
    return createUrl(this.baseUrl, "/api/v1/system/library")
  }
  
  cleanupLibrary() {
    return createUrl(this.baseUrl, "/api/v1/system/library/cleanup")
  }
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/system/library/sse")
  }
  
  updateUserSettings() {
    return createUrl(this.baseUrl, "/api/v1/user/settings")
  }
  
  createApiToken() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
  }
  
  getAllApiTokens() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
  }
  
  deleteApiToken(id: string) {
    return createUrl(this.baseUrl, `/api/v1/user/apitoken/${id}`)
  }
  
  getMedia() {
    return createUrl(this.baseUrl, "/api/v1/media")
  }
  
  getMediaById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}`)
  }
  
  getMediaParts(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts`)
  }
  
  setMediaUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/user`)
  }
  
  getCollections() {
    return createUrl(this.baseUrl, "/api/v1/collections")
  }
  
  getCollectionById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}`)
  }
  
  getCollectionItems(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items`)
  }
  
  createCollection() {
    return createUrl(this.baseUrl, "/api/v1/collections")
  }
  
  editCollection(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}`)
  }
  
  providerMyAnimeListGetAnime(id: string) {
    return createUrl(this.baseUrl, `/api/v1/provider/myanimelist/anime/${id}`)
  }
  
  getMediaImage(id: string, image: string) {
    return createUrl(this.baseUrl, `/files/media/${id}/${image}`)
  }
  
  getCollectionImage(id: string, image: string) {
    return createUrl(this.baseUrl, `/files/collections/${id}/${image}`)
  }
}
