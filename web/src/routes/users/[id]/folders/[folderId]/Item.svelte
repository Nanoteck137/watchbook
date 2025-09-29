<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { MediaType } from "$lib/api-types";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import { getYear } from "$lib/utils";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import {
    ArrowDown,
    ArrowDownToLine,
    ArrowUp,
    ArrowUpToLine,
    EllipsisVertical,
    Trash,
  } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  type Props = {
    folderId: string;
    mediaId: string;

    position: number;
    index: number;
    numItems: number;

    title: string;
    type: MediaType;
    coverUrl?: string | null;
    startDate?: string;

    isUser?: boolean;
  };

  const {
    folderId,
    mediaId,
    position,
    index,
    numItems,
    title,
    type,
    coverUrl,
    startDate,
    isUser,
  }: Props = $props();
  const apiClient = getApiClient();

  const isFirst = $derived(index <= 0);
  const isLast = $derived(index >= numItems - 1);

  const year = $derived(getYear(startDate));
  let openRemoveModal = $state(false);

  async function moveItem(pos: number) {
    const res = await apiClient.moveFolderItem(
      folderId,
      mediaId,
      pos.toString(),
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully moved item");
    invalidateAll();
  }
</script>

<div class="flex items-center justify-between gap-2 border-b py-2">
  <div class="flex items-center gap-4">
    <img
      src={coverUrl}
      alt="Cover"
      class="aspect-[75/106] w-14 rounded object-cover"
    />
    <div>
      <div class="line-clamp-2 text-sm font-medium">{title}</div>
      <div class="text-xs text-muted-foreground">
        {#if year}
          {year} ·
        {/if}
        {type}
      </div>
    </div>
  </div>
  {#if isUser}
    <div class="flex gap-2">
      <Button
        variant="ghost"
        size="icon"
        onclick={() => moveItem(position - 1)}
        disabled={isFirst}
      >
        <ArrowUp />
      </Button>

      <Button
        variant="ghost"
        size="icon"
        onclick={() => moveItem(position + 1)}
        disabled={isLast}
      >
        <ArrowDown />
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={buttonVariants({ variant: "outline", size: "icon" })}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content class="w-40" align="end">
          <DropdownMenu.Group>
            <DropdownMenu.Item onclick={() => moveItem(0)} disabled={isFirst}>
              <ArrowUpToLine />
              Move First
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onclick={() => moveItem(numItems)}
              disabled={isLast}
            >
              <ArrowDownToLine />
              Move Last
            </DropdownMenu.Item>

            <DropdownMenu.Separator />

            <DropdownMenu.Item
              onclick={() => {
                openRemoveModal = true;
              }}
            >
              <Trash />
              Remove
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  {/if}
</div>

<ConfirmBox
  bind:open={openRemoveModal}
  title="Remove Media from Folder?"
  description="Are you sure you want to remove this media item? This action cannot be undone."
  confirmText="Remove Item"
  onResult={async () => {
    const res = await apiClient.removeFolderItem(folderId, mediaId);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully removed folder item");
    invalidateAll();
  }}
/>
