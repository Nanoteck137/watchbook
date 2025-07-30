<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import { onDestroy, onMount } from "svelte";
  import { z } from "zod";

  const { data } = $props();
  const apiClient = getApiClient();

  let test = $state<string[]>([]);
  let syncing = $state(false);

  const SyncError = z.object({
    type: z.string(),
    message: z.string(),
    fullMessage: z.string().optional(),
  });

  const MissingMedia = z.object({
    id: z.string(),
    title: z.string(),
  });

  const Event = z.discriminatedUnion("type", [
    z.object({
      type: z.literal("syncing"),
      data: z.object({
        syncing: z.boolean(),
      }),
    }),
    z.object({
      type: z.literal("report"),
      data: z.object({
        syncErrors: z.array(SyncError).nullable(),
        missingMedia: z.array(MissingMedia).nullable(),
      }),
    }),
  ]);

  onMount(() => {
    console.log("Mount");
    const eventSource = new EventSource(
      data.apiAddress + "/api/v1/system/library/sse",
    );

    eventSource.onmessage = (e) => {
      const event = Event.parse(JSON.parse(e.data));
      console.log(event);

      switch (event.type) {
        case "syncing":
          syncing = event.data.syncing;
          break;
        case "report":
          console.log("Report", event.data);
          // const mapped =
          //   event.data.reports?.map((t) => {
          //     if (t.fullMessage) return t.fullMessage;
          //     return t.message;
          //   }) ?? [];
          // test = mapped;
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

<p>Library Syncing: {syncing}</p>

<Button
  onclick={async () => {
    const res = await apiClient.syncLibrary();
    if (!res.success) {
      handleApiError(res.error);
      return;
    }
  }}
>
  Sync Library
</Button>

<Button
  onclick={async () => {
    const res = await apiClient.cleanupLibrary();
    if (!res.success) {
      handleApiError(res.error);
      return;
    }
  }}
>
  Cleanup Library
</Button>

{#each test as message}
  <p class="whitespace-pre font-mono">{message}</p>
  <br />
{/each}
