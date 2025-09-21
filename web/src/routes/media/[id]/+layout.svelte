<script lang="ts">
  import BannerHeader from "$lib/components/BannerHeader.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Breadcrumb, Button } from "@nanoteck137/nano-ui";
  import Badge from "$lib/components/Badge.svelte";
  import { type UserList } from "$lib/types";
  import SetListModal from "./SetListModal.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import { invalidateAll } from "$app/navigation";
  import toast from "svelte-5-french-toast";
  import { isRoleAdmin } from "$lib/utils";

  const { data, children } = $props();
  const apiClient = getApiClient();

  let openSetListModal = $state(false);
</script>

<Breadcrumb.Root class="py-2">
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/media">Media</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis">
        {data.media.title}
      </Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<Spacer size="md" />

<BannerHeader
  title={data.media.title}
  description={data.media.description}
  coverUrl={data.media.coverUrl}
  bannerUrl={data.media.bannerUrl}
  logoUrl={data.media.logoUrl}
>
  {#snippet buttons()}
    <div class="flex gap-4">
      <Button
        href="/media/{data.media.id}"
        variant="secondary"
        data-sveltekit-noscroll
      >
        Overview
      </Button>
      <Button
        href="/media/{data.media.id}/parts"
        variant="secondary"
        data-sveltekit-noscroll
      >
        Parts
      </Button>
      {#if isRoleAdmin(data.user?.role)}
        <Button
          href="/media/{data.media.id}/settings"
          variant="secondary"
          data-sveltekit-noscroll
        >
          Settings
        </Button>
      {/if}

      <!-- <a class="bg-red-200" href="?">Overview</a>
        <a href="?">Parts</a> -->
      <!-- <Button variant="link">Overview</Button> -->
    </div>
  {/snippet}
  {#snippet underText()}
    <div class="flex gap-2">
      {#if data.user}
        {#if data.media.user?.list}
          <Badge
            class="w-fit hover:cursor-pointer hover:brightness-75"
            list={data.media.user?.list as UserList}
            onclick={() => {
              openSetListModal = true;
            }}
          />

          <p
            class="inline-block w-fit select-none rounded-full bg-blue-500 px-3 py-1 text-xs font-semibold text-white"
          >
            {data.media.user?.currentPart ?? "??"} / {data.media.partCount ??
              "??"}
          </p>
        {:else}
          <Button
            variant="link"
            onclick={async () => {
              const res = await apiClient.setMediaUserData(data.media.id, {
                list: "backlog",
              });
              if (!res.success) {
                return handleApiError(res.error);
              }

              toast.success("Added to user list");
              invalidateAll();
            }}
          >
            Add to list
          </Button>
        {/if}
      {/if}
    </div>
  {/snippet}
</BannerHeader>

<Spacer size="md" />

{@render children()}

<SetListModal
  bind:open={openSetListModal}
  mediaId={data.media.id}
  userList={data.media.user ?? undefined}
/>
