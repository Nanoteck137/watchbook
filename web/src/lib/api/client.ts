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
  
  getUser(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${id}`, "GET", api.GetUser, z.any(), undefined, options)
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
  
  createMedia(body: api.CreateMediaBody, options?: ExtraOptions) {
    return this.request("/api/v1/media", "POST", api.CreateMedia, z.any(), body, options)
  }
  
  editMedia(id: string, body: api.EditMediaBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  changeMediaImages(id: string, body: FormData, options?: ExtraOptions) {
    return this.requestForm(`/api/v1/media/${id}/images`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  addPart(id: string, body: api.AddPartBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/single/parts`, "POST", api.AddPart, z.any(), body, options)
  }
  
  editPart(id: string, index: string, body: api.EditPartBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts/${index}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  removePart(id: string, index: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts/${index}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  setParts(id: string, body: api.SetPartsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts`, "POST", z.undefined(), z.any(), body, options)
  }
  
  setMediaUserData(id: string, body: api.SetMediaUserData, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/user`, "POST", z.undefined(), z.any(), body, options)
  }
  
  deleteMediaUserData(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/user`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  setMediaRelease(id: string, body: api.SetMediaReleaseBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/release`, "POST", z.undefined(), z.any(), body, options)
  }
  
  deleteMediaRelease(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/release`, "DELETE", z.undefined(), z.any(), undefined, options)
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
  
  deleteCollection(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  changeCollectionImages(id: string, body: FormData, options?: ExtraOptions) {
    return this.requestForm(`/api/v1/collections/${id}/images`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  addCollectionItem(id: string, body: api.AddCollectionItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  removeCollectionItem(id: string, mediaId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items/${mediaId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editCollectionItem(id: string, mediaId: string, body: api.EditCollectionItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items/${mediaId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getProviders(options?: ExtraOptions) {
    return this.request("/api/v1/providers", "GET", api.GetProviders, z.any(), undefined, options)
  }
  
  providerSearchMedia(providerName: string, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}`, "GET", api.GetProviderSearch, z.any(), undefined, options)
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
  
  getUser(id: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${id}`)
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
  
  createMedia() {
    return createUrl(this.baseUrl, "/api/v1/media")
  }
  
  editMedia(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}`)
  }
  
  changeMediaImages(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/images`)
  }
  
  addPart(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/single/parts`)
  }
  
  editPart(id: string, index: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts/${index}`)
  }
  
  removePart(id: string, index: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts/${index}`)
  }
  
  setParts(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts`)
  }
  
  setMediaUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/user`)
  }
  
  deleteMediaUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/user`)
  }
  
  setMediaRelease(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/release`)
  }
  
  deleteMediaRelease(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/release`)
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
  
  deleteCollection(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}`)
  }
  
  changeCollectionImages(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/images`)
  }
  
  addCollectionItem(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items`)
  }
  
  removeCollectionItem(id: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items/${mediaId}`)
  }
  
  editCollectionItem(id: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items/${mediaId}`)
  }
  
  getProviders() {
    return createUrl(this.baseUrl, "/api/v1/providers")
  }
  
  providerSearchMedia(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}`)
  }
  
  getMediaImage(id: string, file: string) {
    return createUrl(this.baseUrl, `/files/media/${id}/images/${file}`)
  }
  
  getCollectionImage(id: string, file: string) {
    return createUrl(this.baseUrl, `/files/collections/${id}/images/${file}`)
  }
}
