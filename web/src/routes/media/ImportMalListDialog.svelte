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
  <Dialog.Trigger class={cn(buttonVariants())}>Import MAL List</Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <form
      class="w-full"
      onsubmit={async (e) => {
        e.preventDefault();

        const formData = new FormData(e.target as HTMLFormElement);

        const username = formData.get("username")?.toString() ?? "0";
        const overrideExistingEntries =
          formData.get("overrideExistingEntries")?.toString() ?? "";

        const res = await apiClient.importMalList({
          username: username,
          overrideExistingEntries: overrideExistingEntries === "on",
        });
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully imported user list");
        invalidateAll();
        open = false;
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Import MyAnimeList user list</Dialog.Title>
        <!-- <Dialog.Description>
          Make changes to your profile here. Click save when you're done.
        </Dialog.Description> -->
      </Dialog.Header>
      <div class="grid gap-4 py-4">
        <div class="flex flex-col gap-2">
          <Label for="username">MyAnimeList username</Label>
          <Input name="username" id="username" />
        </div>
        <div class="flex items-center gap-2">
          <Checkbox
            name="overrideExistingEntries"
            id="overrideExistingEntries"
            checked={false}
          />
          <Label for="overrideExistingEntries">Override Existing Entries</Label
          >
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
