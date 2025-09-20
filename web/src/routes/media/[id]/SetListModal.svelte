<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import {
    Button,
    Checkbox,
    Dialog,
    Input,
    Label,
    Select,
  } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { MediaUser } from "$lib/api/types";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import { Trash } from "lucide-svelte";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import {
    MediaUserListEnum,
    type MediaUserList,
    mediaUserLists,
  } from "$lib/api-types";

  const Schema = z.object({
    list: MediaUserListEnum,
    score: z.string(),
    currentPart: z.number().min(0),
    revisitCount: z.number().min(0),
    isRevisiting: z.boolean(),
  });
  type SchemaTy = z.infer<typeof Schema>;

  export type Props = {
    open: boolean;
    mediaId: string;
    userList?: MediaUser;
  };

  let { open = $bindable(), mediaId, userList }: Props = $props();
  const apiClient = getApiClient();

  let openDeleteConfirmModal = $state(false);

  $effect(() => {
    if (open) {
      reset({
        data: {
          list: (userList?.list as MediaUserList) ?? "backlog",
          score: userList?.score?.toString() ?? "0",
          currentPart: userList?.currentPart ?? 0,
          revisitCount: userList?.revisitCount ?? 0,
          isRevisiting: userList?.isRevisiting ?? false,
        },
      });
    }
  });

  async function submit(data: SchemaTy) {
    const res = await apiClient.setMediaUserData(mediaId, {
      list: data.list,
      score: parseInt(data.score),
      currentPart: data.currentPart,
      revisitCount: data.revisitCount,
      isRevisiting: data.isRevisiting,
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    open = false;

    toast.success("Successfully updated list");
    invalidateAll();
  }

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          await submit(form.data);
        }
      },
    },
  );

  const scores = [
    { label: "(10) Masterpiece", value: "10" },
    { label: "(9) Great", value: "9" },
    { label: "(8) Very Good", value: "8" },
    { label: "(7) Good", value: "7" },
    { label: "(6) Fine", value: "6" },
    { label: "(5) Average", value: "5" },
    { label: "(4) Bad", value: "4" },
    { label: "(3) Very Bad", value: "3" },
    { label: "(2) Horrible", value: "2" },
    { label: "(1) Appalling", value: "1" },
    { label: "No Score", value: "0" },
  ] as const;
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Set list</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="list">List</Label>

        <div class="flex items-center gap-2">
          <Select.Root
            type="single"
            bind:value={$form.list}
            allowDeselect={false}
          >
            <Select.Trigger>
              {mediaUserLists.find((f) => f.value === $form.list)?.label ??
                "Select list"}
            </Select.Trigger>
            <Select.Content>
              {#each mediaUserLists as list (list.value)}
                <Select.Item value={list.value} label={list.label} />
              {/each}
            </Select.Content>
          </Select.Root>

          <Button
            variant="destructive"
            onclick={() => {
              openDeleteConfirmModal = true;
            }}
          >
            <Trash />
          </Button>
        </div>

        <Errors errors={$errors.list} />
      </FormItem>

      <FormItem>
        <Label for="score">Score</Label>

        <Select.Root
          type="single"
          bind:value={$form.score}
          allowDeselect={false}
        >
          <Select.Trigger>
            {scores.find((f) => f.value === $form.score)?.label ??
              "Select score"}
          </Select.Trigger>
          <Select.Content>
            {#each scores as score (score.value)}
              <Select.Item value={score.value} label={score.label} />
            {/each}
          </Select.Content>
        </Select.Root>
      </FormItem>

      <FormItem>
        <Label for="currentPart">Current Part</Label>
        <Input
          id="currentPart"
          name="currentPart"
          type="number"
          min={0}
          bind:value={$form.currentPart}
        />
        <Errors errors={$errors.currentPart} />
      </FormItem>

      <FormItem>
        <Label for="revistCount">Revisting Count</Label>
        <Input
          id="revistCount"
          name="revistCount"
          type="number"
          min={0}
          bind:value={$form.revisitCount}
        />
        <Errors errors={$errors.revisitCount} />
      </FormItem>

      <FormItem>
        <div class="flex items-center gap-2">
          <Checkbox
            id="isRevisiting"
            name="isRevisiting"
            bind:checked={$form.isRevisiting}
          />
          <Label for="isRevisiting">Is Revisiting</Label>
        </div>
        <Errors errors={$errors.isRevisiting} />
      </FormItem>

      <!-- <FormItem>
        <Label for="episode">Current episode</Label>
        <Input name="episode" id="episode" type="number" />
      </FormItem> -->

      <!-- <div class="flex flex-col gap-2">
        <Label for="rewatchCount">Rewatch Count</Label>
        <Input
          name="rewatchCount"
          id="rewatchCount"
          type="number"
          value={data.media.user?.revisitCount ?? 0}
        />
      </div>
      <div class="flex items-center gap-2">
        <Checkbox
          name="isRewatching"
          id="isRewatching"
          checked={data.media.user?.isRevisiting ?? false}
        />
        <Label for="isRewatching">Is Rewatching</Label>
      </div> -->

      <Dialog.Footer class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          onclick={() => {
            open = false;
            reset();
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Update
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>

<ConfirmBox
  bind:open={openDeleteConfirmModal}
  title="Remove from list?"
  description="Are you sure you want to remove this from the list? This action cannot be undone."
  onResult={async () => {
    const res = await apiClient.deleteMediaUserData(mediaId);
    if (!res.success) {
      return handleApiError(res.error);
    }

    open = false;
    invalidateAll();
  }}
/>
