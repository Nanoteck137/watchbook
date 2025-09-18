<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import { cn } from "$lib/utils";
  import { buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import {
    Download,
    EllipsisVertical,
    Image,
    Pencil,
    Trash,
  } from "lucide-svelte";
  import type { Media } from "$lib/api/types";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import EditImagesModal from "./EditImagesModal.svelte";
  import toast from "svelte-5-french-toast";
  import { goto } from "$app/navigation";
  import EditMediaModal from "./EditMediaModal.svelte";
  import ProviderUpdateModal from "./ProviderUpdateModal.svelte";

  type Props = {
    media: Media;
  };

  const { media }: Props = $props();

  const apiClient = getApiClient();

  let openEditModal = $state(false);
  let openEditImagesModal = $state(false);
  let openDeleteModal = $state(false);
  let openProviderUpdateModal = $state(false);
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger
    class={cn(
      "absolute right-2 top-2 z-20",
      buttonVariants({ variant: "secondary", size: "icon" }),
    )}
  >
    <EllipsisVertical />
  </DropdownMenu.Trigger>
  <DropdownMenu.Content class="w-40" align="center">
    <DropdownMenu.Group>
      <DropdownMenu.Item
        onclick={() => {
          openProviderUpdateModal = true;
        }}
      >
        <Download />
        Update
      </DropdownMenu.Item>

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
          openEditImagesModal = true;
        }}
      >
        <Image />
        Edit Images
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onclick={() => {
          openDeleteModal = true;
        }}
      >
        <Trash />
        Delete
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>

<EditMediaModal bind:open={openEditModal} {media} />

<EditImagesModal bind:open={openEditImagesModal} mediaId={media.id} />

<ConfirmBox
  bind:open={openDeleteModal}
  title="Delete Media?"
  description="Are you sure you want to delete this media? This action cannot be undone."
  confirmText="Delete"
  onResult={async () => {
    const res = await apiClient.deleteMedia(media.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted media");
    goto(`/media`, { invalidateAll: true });
  }}
/>

<ProviderUpdateModal
  bind:open={openProviderUpdateModal}
  mediaId={media.id}
  providers={media.providers}
/>
