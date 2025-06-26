<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { cn } from "$lib/utils";
  import {
    Button,
    buttonVariants,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  let open = $state(false);
  const apiClient = getApiClient();
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger class={cn(buttonVariants())}>
    Import Anime from MAL
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <form
      class="w-full"
      onsubmit={async (e) => {
        e.preventDefault();

        const formData = new FormData(e.target as HTMLFormElement);

        const entry = formData.get("entry")?.toString() ?? "0";

        // TODO(patrik): Handle errors
        let id = parseInt(entry);
        if (isNaN(id)) {
          const url = new URL(entry);
          if (url.host === "myanimelist.net") {
            const splits = url.pathname.split("/");
            if (splits.length >= 3) {
              if (splits[1] === "anime") {
                id = parseInt(splits[2]);
              }
            }
          }
        }

        if (isNaN(id)) {
          // ERROR
          return;
        }

        const res = await apiClient.importMalAnime({
          id: id.toString(),
        });
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully added anime");
        open = false;

        goto(`/animes/${res.data.animeId}`, { invalidateAll: true });
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Import MyAnimeList user list</Dialog.Title>
      </Dialog.Header>
      <div class="grid gap-4 py-4">
        <div class="flex flex-col gap-2">
          <Label for="entry">MyAnimeList url/id</Label>
          <Input name="entry" id="entry" />
        </div>
      </div>
      <Dialog.Footer>
        <Button
          variant="ghost"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button type="submit">Save changes</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
