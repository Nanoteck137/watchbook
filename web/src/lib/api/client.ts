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
  
  startDownload(body: api.StartDownloadBody, options?: ExtraOptions) {
    return this.request("/api/v1/system/download", "POST", z.undefined(), z.any(), body, options)
  }
  
  cancelDownload(options?: ExtraOptions) {
    return this.request("/api/v1/system/download", "DELETE", z.undefined(), z.any(), undefined, options)
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
  
  createMedia(body: api.CreateMediaBody, options?: ExtraOptions) {
    return this.request("/api/v1/media", "POST", api.CreateMedia, z.any(), body, options)
  }
  
  editMedia(id: string, body: api.EditMediaBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getMediaParts(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts`, "GET", api.GetMediaParts, z.any(), undefined, options)
  }
  
  addPart(id: string, body: api.AddPartBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/single/parts`, "POST", api.AddPart, z.any(), body, options)
  }
  
  addMultipleParts(id: string, body: api.AddMultiplePartsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/multiple/parts`, "POST", z.undefined(), z.any(), body, options)
  }
  
  editPart(id: string, index: string, body: api.EditPartBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts/${index}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  removePart(id: string, index: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts/${index}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  addImage(id: string, body: api.AddImageBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/images`, "POST", api.AddImage, z.any(), body, options)
  }
  
  editImage(id: string, hash: string, body: api.EditImageBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/images/${hash}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  setMediaUserData(id: string, body: api.SetMediaUserData, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/user`, "POST", z.undefined(), z.any(), body, options)
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
  
  startDownload() {
    return createUrl(this.baseUrl, "/api/v1/system/download")
  }
  
  cancelDownload() {
    return createUrl(this.baseUrl, "/api/v1/system/download")
  }
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/system/sse")
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
  
  createMedia() {
    return createUrl(this.baseUrl, "/api/v1/media")
  }
  
  editMedia(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}`)
  }
  
  getMediaParts(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts`)
  }
  
  addPart(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/single/parts`)
  }
  
  addMultipleParts(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/multiple/parts`)
  }
  
  editPart(id: string, index: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts/${index}`)
  }
  
  removePart(id: string, index: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts/${index}`)
  }
  
  addImage(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/images`)
  }
  
  editImage(id: string, hash: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/images/${hash}`)
  }
  
  setMediaUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/user`)
  }
  
  providerMyAnimeListGetAnime(id: string) {
    return createUrl(this.baseUrl, `/api/v1/provider/myanimelist/anime/${id}`)
  }
  
  getMediaImage(id: string, image: string) {
    return createUrl(this.baseUrl, `/files/media/${id}/${image}`)
  }
}
