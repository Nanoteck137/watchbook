<script lang="ts">
  import MediaCard from "$lib/components/MediaCard.svelte";
  import { Breadcrumb, Button, buttonVariants } from "@nanoteck137/nano-ui";
  import { FileQuestion, Image as ImageIcon, Plus } from "lucide-svelte";
  import ShowLogoModal from "./ShowLogoModal.svelte";
  import { cn, isRoleAdmin } from "$lib/utils";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import MediaItemDropdown from "./MediaItemDropdown.svelte";
  import AddMediaItem from "./AddMediaItem.svelte";
  import EditImagesModal from "./EditImagesModal.svelte";
  import CollectionDropdown from "./CollectionDropdown.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import BannerHeader from "$lib/components/BannerHeader.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openAddMediaModal = $state(false);
  let openEditImagesModal = $state(false);
</script>

<Breadcrumb.Root class="py-2">
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/collections">Collections</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis">
        {data.collection.name}
      </Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<Spacer size="md" />

<BannerHeader
  title={data.collection.name}
  coverUrl={data.collection.coverUrl}
  bannerUrl={data.collection.bannerUrl}
  logoUrl={data.collection.logoUrl}
>
  {#snippet imageContent()}
    <CollectionDropdown collection={data.collection} />
  {/snippet}
</BannerHeader>

<Spacer size="md" />

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    Collection Items
    {#if isRoleAdmin(data.user?.role)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openAddMediaModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </h2>
  <p class="text-sm">{data.items.length} item(s)</p>
</div>

<Spacer size="md" />

<div
  class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
>
  {#each data.items as item}
    <div class="relative">
      <MediaCard
        href="/media/{item.mediaId}"
        coverUrl={item.coverUrl}
        title={item.collectionName}
        startDate={item.startDate}
        partCount={item.partCount}
        score={item.score}
        userList={item.user?.list ?? null}
      />

      <MediaItemDropdown {item} />
    </div>
  {/each}
</div>

<AddMediaItem
  bind:open={openAddMediaModal}
  itemIds={data.items.map((i) => i.mediaId)}
  onResult={async (results) => {
    for (const item of results) {
      const res = await apiClient.addCollectionItem(data.collection.id, {
        mediaId: item.id,
        name: item.name,
        searchSlug: item.name,
        position: item.position,
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    toast.success("Successfully added new media items");
    invalidateAll();
  }}
/>

<EditImagesModal
  bind:open={openEditImagesModal}
  collectionId={data.collection.id}
/>
