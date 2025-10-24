<script lang="ts">
  import type { ShowSeason } from "$lib/api/types";
  import { cn, isRoleAdmin } from "$lib/utils";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import AddMediaItem from "./AddMediaItem.svelte";
  import {
    EllipsisVertical,
    Pencil,
    Plus,
    Settings,
    Trash,
  } from "lucide-svelte";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import MediaCard from "$lib/components/MediaCard.svelte";
  import MediaItemDropdown from "./MediaItemDropdown.svelte";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import EditSeasonModal from "./EditSeasonModal.svelte";

  type Props = {
    userRole?: string;
    season: ShowSeason;
  };

  const { userRole, season }: Props = $props();
  const apiClient = getApiClient();

  let openAddMediaItemModal = $state(false);
  let openEditModal = $state(false);
  let openRemoveModal = $state(false);
</script>

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    {season.num} - {season.name}
    {#if isRoleAdmin(userRole)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openAddMediaItemModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}

    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class={cn(buttonVariants({ variant: "ghost", size: "icon" }))}
      >
        <Settings />
      </DropdownMenu.Trigger>
      <DropdownMenu.Content class="w-40" align="start">
        <DropdownMenu.Group>
          <DropdownMenu.Item
            onclick={() => {
              openEditModal = true;
            }}
          >
            <Pencil />
            Edit
          </DropdownMenu.Item>

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
  </h2>
  <!-- <p class="text-sm">{data.items.length} item(s)</p> -->
</div>

<div class="overflow-x flex gap-2 p-6">
  {#each season.items as item}
    <div class="relative">
      <MediaCard
        href="/media/{item.mediaId}"
        coverUrl={item.coverUrl}
        title={item.title}
        partCount={item.partCount}
        score={item.score}
        userList={item.user?.list ?? null}
      />

      <MediaItemDropdown {item} />
    </div>
  {/each}
</div>

<AddMediaItem
  bind:open={openAddMediaItemModal}
  itemIds={season.items.map((i) => i.mediaId)}
  onResult={async (items) => {
    for (const item of items) {
      const res = await apiClient.addShowSeasonItem(
        season.showId,
        season.num.toString(),
        {
          mediaId: item.id,
          position: item.position,
        },
      );
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    toast.success("Successfully added new media items");
    invalidateAll();
  }}
/>

<EditSeasonModal bind:open={openEditModal} {season} />

<ConfirmBox
  bind:open={openRemoveModal}
  title="Remove Season?"
  description="Are you sure you want to remove this season? This action cannot be undone."
  confirmText="Remove Season"
  onResult={async () => {
    const res = await apiClient.removeShowSeason(
      season.showId,
      season.num.toString(),
    );
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully removed season");
    invalidateAll();
  }}
/>
