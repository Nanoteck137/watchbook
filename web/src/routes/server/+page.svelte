<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import { onMount } from "svelte";
  import toast from "svelte-5-french-toast";
  import { z } from "zod";

  const { data } = $props();
  const apiClient = getApiClient();

  let isDownloading = $state(false);
  let lastError = $state("");
  let currentDownload = $state(0);
  let totalDownloads = $state(0);

  const Event = z.discriminatedUnion("type", [
    z.object({
      type: z.literal("status"),
      data: z.object({
        isDownloading: z.boolean(),
        currentDownload: z.number(),
        totalDownloads: z.number(),
        lastError: z.string(),
      }),
    }),
  ]);

  onMount(() => {
    console.log("Mount");
    const eventSource = new EventSource(
      data.apiAddress + "/api/v1/system/sse",
    );

    eventSource.onmessage = (e) => {
      const event = Event.parse(JSON.parse(e.data));
      console.log(event);

      switch (event.type) {
        case "status":
          isDownloading = event.data.isDownloading;
          currentDownload = event.data.currentDownload;
          totalDownloads = event.data.totalDownloads;
          lastError = event.data.lastError;
          break;
      }
    };

    return () => {
      console.log("Cleanup");
      eventSource.close();
    };
  });
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

<Button
  onclick={async () => {
    // const filter = "lastDataFetch == null";
    const filter = 'status != "finished" || lastDataFetch == null';

    const getAnime = async function (page: number) {
      const animes = await apiClient.getAnimes({
        query: { filter, page: page.toString() },
      });
      if (!animes.success) {
        handleApiError(animes.error);
        throw "Failed";
      }

      return animes.data;
    };

    const firstPage = await getAnime(0);
    let ids = firstPage.animes.map((a) => a.id);

    for (let i = 1; i < firstPage.page.totalPages; i++) {
      const anime = await getAnime(i);
      anime.animes.forEach((a) => {
        ids.push(a.id);
      });
    }

    const res = await apiClient.startDownload({ ids });
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Starting download");
  }}
>
  Download Not Fetched Animes
</Button>

<Button
  onclick={async () => {
    // const res = await apiClient.startDownload();
    // if (!res.success) {
    //   return handleApiError(res.error);
    // }
    // toast.success("Starting download");
  }}>Start Download</Button
>

<p>Is Downloading: {isDownloading}</p>
{#if lastError.length > 0}
  <p>Last Error: {lastError}</p>
{/if}
{#if isDownloading}
  <p>
    Progress: {currentDownload} / {totalDownloads} ({Math.floor(
      (currentDownload / totalDownloads) * 100,
    )}%)
  </p>
  <Button
    onclick={async () => {
      const res = await apiClient.cancelDownload();
      if (!res.success) {
        return handleApiError(res.error);
      }
      toast.success("Canceled download");
    }}
  >
    Cancel Download
  </Button>
{/if}
