import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, createUrl, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  url: ClientUrls;

  constructor(baseUrl: string) {
    super(baseUrl);
    this.url = new ClientUrls(baseUrl);
  }
  
  addCollectionItem(id: string, body: api.AddCollectionItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  addFolderItem(id: string, mediaId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}/items/${mediaId}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  addPart(id: string, body: api.AddPartBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/single/parts`, "POST", api.AddPart, z.any(), body, options)
  }
  
  addShowSeason(id: string, body: api.AddShowSeasonBody, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons`, "POST", z.undefined(), z.any(), body, options)
  }
  
  addShowSeasonItem(id: string, seasonNum: string, body: api.AddShowSeasonItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  changeCollectionImages(id: string, body: FormData, options?: ExtraOptions) {
    return this.requestForm(`/api/v1/collections/${id}/images`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  changePassword(body: api.ChangePasswordBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/password", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  changeShowImages(id: string, body: FormData, options?: ExtraOptions) {
    return this.requestForm(`/api/v1/shows/${id}/images`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  createCollection(body: api.CreateCollectionBody, options?: ExtraOptions) {
    return this.request("/api/v1/collections", "POST", api.CreateCollection, z.any(), body, options)
  }
  
  createFolder(body: api.CreateFolderBody, options?: ExtraOptions) {
    return this.request("/api/v1/folders", "POST", api.CreateFolder, z.any(), body, options)
  }
  
  createMedia(body: api.CreateMediaBody, options?: ExtraOptions) {
    return this.request("/api/v1/media", "POST", api.CreateMedia, z.any(), body, options)
  }
  
  createShow(body: api.CreateShowBody, options?: ExtraOptions) {
    return this.request("/api/v1/shows", "POST", api.CreateShow, z.any(), body, options)
  }
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteCollection(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteFolder(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteMedia(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteMediaRelease(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/release`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteMediaUserData(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/user`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteShow(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editCollection(id: string, body: api.EditCollectionBody, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editCollectionItem(id: string, mediaId: string, body: api.EditCollectionItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items/${mediaId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editFolder(id: string, body: api.EditFolderBody, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editMedia(id: string, body: api.EditMediaBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editPart(id: string, index: string, body: api.EditPartBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts/${index}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editShow(id: string, body: api.EditShowBody, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editShowSeason(id: string, seasonNum: string, body: api.EditShowSeasonBody, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editShowSeasonItem(id: string, seasonNum: string, mediaId: string, body: api.EditShowSeasonItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}/items/${mediaId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getAllApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "GET", api.GetAllApiTokens, z.any(), undefined, options)
  }
  
  getCollectionById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "GET", api.GetCollectionById, z.any(), undefined, options)
  }
  
  
  getCollectionItems(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items`, "GET", api.GetCollectionItems, z.any(), undefined, options)
  }
  
  getCollections(options?: ExtraOptions) {
    return this.request("/api/v1/collections", "GET", api.GetCollections, z.any(), undefined, options)
  }
  
  getFolderById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}`, "GET", api.GetFolderById, z.any(), undefined, options)
  }
  
  getFolderItems(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}/items`, "GET", api.GetFolderItems, z.any(), undefined, options)
  }
  
  getFolders(options?: ExtraOptions) {
    return this.request("/api/v1/folders", "GET", api.GetFolders, z.any(), undefined, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
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
  
  getProviders(options?: ExtraOptions) {
    return this.request("/api/v1/providers", "GET", api.GetProviders, z.any(), undefined, options)
  }
  
  getShowById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}`, "GET", api.GetShowById, z.any(), undefined, options)
  }
  
  
  getShowSeason(id: string, seasonNum: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}`, "GET", api.GetShowSeason, z.any(), undefined, options)
  }
  
  getShowSeasonEpisodes(id: string, seasonNum: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}/episodes`, "GET", api.GetShowSeasonEpisodes, z.any(), undefined, options)
  }
  
  getShowSeasons(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons`, "GET", api.GetShowSeasons, z.any(), undefined, options)
  }
  
  getShows(options?: ExtraOptions) {
    return this.request("/api/v1/shows", "GET", api.GetShows, z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  getUser(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${id}`, "GET", api.GetUser, z.any(), undefined, options)
  }
  
  getUserStats(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${id}/stats`, "GET", api.GetUserStats, z.any(), undefined, options)
  }
  
  importMalAnimeList(username: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/import/mal/${username}/anime`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  moveFolderItem(id: string, mediaId: string, pos: string, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}/items/${mediaId}/move/${pos}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  providerImportCollections(providerName: string, body: api.PostProviderImportCollectionsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/collections/import`, "POST", z.undefined(), z.any(), body, options)
  }
  
  providerImportMedia(providerName: string, body: api.PostProviderImportMediaBody, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/media/import`, "POST", z.undefined(), z.any(), body, options)
  }
  
  providerImportShows(providerName: string, body: api.PostProviderImportCollectionsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/shows/import`, "POST", z.undefined(), z.any(), body, options)
  }
  
  providerSearchCollections(providerName: string, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/collections`, "GET", api.GetProviderSearch, z.any(), undefined, options)
  }
  
  providerSearchMedia(providerName: string, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/media`, "GET", api.GetProviderSearch, z.any(), undefined, options)
  }
  
  providerSearchShows(providerName: string, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/shows`, "GET", api.GetProviderSearch, z.any(), undefined, options)
  }
  
  providerUpdateCollection(providerName: string, collectionId: string, body: api.ProviderCollectionUpdateBody, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/collections/${collectionId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  providerUpdateMedia(providerName: string, mediaId: string, body: api.ProviderMediaUpdateBody, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/media/${mediaId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  providerUpdateShow(providerName: string, showId: string, body: api.ProviderCollectionUpdateBody, options?: ExtraOptions) {
    return this.request(`/api/v1/providers/${providerName}/collections/${showId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  providerUpdateUnknownMedia(options?: ExtraOptions) {
    return this.request("/api/v1/providers/updateUnknownMedia", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  removeCollectionItem(id: string, mediaId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}/items/${mediaId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  removeFolderItem(id: string, mediaId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/folders/${id}/items/${mediaId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  removePart(id: string, index: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts/${index}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  removeShowSeason(id: string, seasonNum: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  removeShowSeasonItem(id: string, seasonNum: string, mediaId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/shows/${id}/seasons/${seasonNum}/items/${mediaId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  setMediaRelease(id: string, body: api.SetMediaReleaseBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/release`, "POST", z.undefined(), z.any(), body, options)
  }
  
  setMediaUserData(id: string, body: api.SetMediaUserData, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/user`, "POST", z.undefined(), z.any(), body, options)
  }
  
  setParts(id: string, body: api.SetPartsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}/parts`, "POST", z.undefined(), z.any(), body, options)
  }
  
  signin(body: api.SigninBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signin", "POST", api.Signin, z.any(), body, options)
  }
  
  signup(body: api.SignupBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signup", "POST", api.Signup, z.any(), body, options)
  }
  
  test(options?: ExtraOptions) {
    return this.request("/api/v1/test", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  test2(options?: ExtraOptions) {
    return this.request("/api/v1/test2", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
}

export class ClientUrls {
  baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }
  
  addCollectionItem(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items`)
  }
  
  addFolderItem(id: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}/items/${mediaId}`)
  }
  
  addPart(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/single/parts`)
  }
  
  addShowSeason(id: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons`)
  }
  
  addShowSeasonItem(id: string, seasonNum: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}/items`)
  }
  
  changeCollectionImages(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/images`)
  }
  
  changePassword() {
    return createUrl(this.baseUrl, "/api/v1/auth/password")
  }
  
  changeShowImages(id: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/images`)
  }
  
  createApiToken() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
  }
  
  createCollection() {
    return createUrl(this.baseUrl, "/api/v1/collections")
  }
  
  createFolder() {
    return createUrl(this.baseUrl, "/api/v1/folders")
  }
  
  createMedia() {
    return createUrl(this.baseUrl, "/api/v1/media")
  }
  
  createShow() {
    return createUrl(this.baseUrl, "/api/v1/shows")
  }
  
  deleteApiToken(id: string) {
    return createUrl(this.baseUrl, `/api/v1/user/apitoken/${id}`)
  }
  
  deleteCollection(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}`)
  }
  
  deleteFolder(id: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}`)
  }
  
  deleteMedia(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}`)
  }
  
  deleteMediaRelease(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/release`)
  }
  
  deleteMediaUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/user`)
  }
  
  deleteShow(id: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}`)
  }
  
  editCollection(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}`)
  }
  
  editCollectionItem(id: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items/${mediaId}`)
  }
  
  editFolder(id: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}`)
  }
  
  editMedia(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}`)
  }
  
  editPart(id: string, index: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts/${index}`)
  }
  
  editShow(id: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}`)
  }
  
  editShowSeason(id: string, seasonNum: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}`)
  }
  
  editShowSeasonItem(id: string, seasonNum: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}/items/${mediaId}`)
  }
  
  getAllApiTokens() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
  }
  
  getCollectionById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}`)
  }
  
  getCollectionImage(id: string, file: string) {
    return createUrl(this.baseUrl, `/files/collections/${id}/images/${file}`)
  }
  
  getCollectionItems(id: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items`)
  }
  
  getCollections() {
    return createUrl(this.baseUrl, "/api/v1/collections")
  }
  
  getFolderById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}`)
  }
  
  getFolderItems(id: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}/items`)
  }
  
  getFolders() {
    return createUrl(this.baseUrl, "/api/v1/folders")
  }
  
  getMe() {
    return createUrl(this.baseUrl, "/api/v1/auth/me")
  }
  
  getMedia() {
    return createUrl(this.baseUrl, "/api/v1/media")
  }
  
  getMediaById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}`)
  }
  
  getMediaImage(id: string, file: string) {
    return createUrl(this.baseUrl, `/files/media/${id}/images/${file}`)
  }
  
  getMediaParts(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts`)
  }
  
  getProviders() {
    return createUrl(this.baseUrl, "/api/v1/providers")
  }
  
  getShowById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}`)
  }
  
  getShowImage(id: string, file: string) {
    return createUrl(this.baseUrl, `/files/shows/${id}/images/${file}`)
  }
  
  getShowSeason(id: string, seasonNum: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}`)
  }
  
  getShowSeasonEpisodes(id: string, seasonNum: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}/episodes`)
  }
  
  getShowSeasons(id: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons`)
  }
  
  getShows() {
    return createUrl(this.baseUrl, "/api/v1/shows")
  }
  
  getSystemInfo() {
    return createUrl(this.baseUrl, "/api/v1/system/info")
  }
  
  getUser(id: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${id}`)
  }
  
  getUserStats(id: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${id}/stats`)
  }
  
  importMalAnimeList(username: string) {
    return createUrl(this.baseUrl, `/api/v1/users/import/mal/${username}/anime`)
  }
  
  moveFolderItem(id: string, mediaId: string, pos: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}/items/${mediaId}/move/${pos}`)
  }
  
  providerImportCollections(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/collections/import`)
  }
  
  providerImportMedia(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/media/import`)
  }
  
  providerImportShows(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/shows/import`)
  }
  
  providerSearchCollections(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/collections`)
  }
  
  providerSearchMedia(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/media`)
  }
  
  providerSearchShows(providerName: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/shows`)
  }
  
  providerUpdateCollection(providerName: string, collectionId: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/collections/${collectionId}`)
  }
  
  providerUpdateMedia(providerName: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/media/${mediaId}`)
  }
  
  providerUpdateShow(providerName: string, showId: string) {
    return createUrl(this.baseUrl, `/api/v1/providers/${providerName}/collections/${showId}`)
  }
  
  providerUpdateUnknownMedia() {
    return createUrl(this.baseUrl, "/api/v1/providers/updateUnknownMedia")
  }
  
  removeCollectionItem(id: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/collections/${id}/items/${mediaId}`)
  }
  
  removeFolderItem(id: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/folders/${id}/items/${mediaId}`)
  }
  
  removePart(id: string, index: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts/${index}`)
  }
  
  removeShowSeason(id: string, seasonNum: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}`)
  }
  
  removeShowSeasonItem(id: string, seasonNum: string, mediaId: string) {
    return createUrl(this.baseUrl, `/api/v1/shows/${id}/seasons/${seasonNum}/items/${mediaId}`)
  }
  
  setMediaRelease(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/release`)
  }
  
  setMediaUserData(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/user`)
  }
  
  setParts(id: string) {
    return createUrl(this.baseUrl, `/api/v1/media/${id}/parts`)
  }
  
  signin() {
    return createUrl(this.baseUrl, "/api/v1/auth/signin")
  }
  
  signup() {
    return createUrl(this.baseUrl, "/api/v1/auth/signup")
  }
  
  test() {
    return createUrl(this.baseUrl, "/api/v1/test")
  }
  
  test2() {
    return createUrl(this.baseUrl, "/api/v1/test2")
  }
  
  updateUserSettings() {
    return createUrl(this.baseUrl, "/api/v1/user/settings")
  }
}
