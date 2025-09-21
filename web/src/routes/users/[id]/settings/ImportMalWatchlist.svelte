<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";

  const Schema = z.object({
    username: z.string().min(1),
  });

  const apiClient = getApiClient();

  const f = superForm(defaults(zod(Schema)), {
    id: "import-mal-watchlist",
    SPA: true,
    validators: zod(Schema),
    dataType: "json",
    async onUpdate({ form }) {
      if (form.valid) {
        const data = form.data;

        const res = await apiClient.importMalAnimeList(data.username);
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Submitted import job request");
        invalidateAll();
      }
    },
  });
  const { form, errors, enhance, submitting } = f;
</script>

<div class="flex flex-col items-center gap-4 border-b p-6">
  <h2 class="text-bold text-center text-xl">Import MyAnimeList Watchlist</h2>

  <form class="flex w-full flex-col gap-4 sm:max-w-[460px]" use:enhance>
    <FormItem>
      <Label for="username">Username</Label>
      <Input
        id="username"
        name="username"
        type="text"
        bind:value={$form.username}
      />
      <Errors errors={$errors.username} />
    </FormItem>

    <div class="flex flex-col justify-end sm:flex-row">
      <Button type="submit" disabled={$submitting}>
        Send Import Request
        {#if $submitting}
          <Spinner />
        {/if}
      </Button>
    </div>
  </form>
</div>
