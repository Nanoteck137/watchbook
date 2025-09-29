<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, setError, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";

  const Schema = z.object({
    currentPassword: z.string().min(1),
    newPassword: z.string().min(1),
    confirmPassword: z.string().min(1),
  });

  const apiClient = getApiClient();

  const { form, errors, enhance, submitting } = superForm(
    defaults(zod(Schema)),
    {
      id: "change-password",
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;
          const res = await apiClient.changePassword({
            currentPassword: formData.currentPassword,
            newPassword: formData.newPassword,
            newPasswordConfirm: formData.confirmPassword,
          });
          if (!res.success) {
            if (res.error.type === "VALIDATION_ERROR") {
              setError(form, "newPassword", res.error.extra.newPassword);
              setError(
                form,
                "confirmPassword",
                res.error.extra.newPasswordConfirm,
              );
            }

            if (res.error.type === "INVALID_CREDENTIALS") {
              setError(form, "currentPassword", res.error.message);
            }

            return handleApiError(res.error);
          }

          toast.success("Successfully changed password");
          invalidateAll();
        }
      },
    },
  );
</script>

<div class="flex flex-col items-center gap-4 border-b p-6">
  <h2 class="text-bold text-center text-xl">Change Password</h2>

  <form class="flex w-full flex-col gap-4 sm:max-w-[460px]" use:enhance>
    <FormItem>
      <Label for="currentPassword">Current Password</Label>
      <Input
        id="currentPassword"
        name="currentPassword"
        type="password"
        bind:value={$form.currentPassword}
      />
      <Errors errors={$errors.currentPassword} />
    </FormItem>

    <FormItem>
      <Label for="newPassword">New Password</Label>
      <Input
        id="newPassword"
        name="newPassword"
        type="password"
        bind:value={$form.newPassword}
      />
      <Errors errors={$errors.newPassword} />
    </FormItem>

    <FormItem>
      <Label for="confirmPassword">Confirm New Password</Label>
      <Input
        id="confirmPassword"
        name="confirmPassword"
        type="password"
        bind:value={$form.confirmPassword}
      />
      <Errors errors={$errors.confirmPassword} />
    </FormItem>

    <div class="flex flex-col justify-end sm:flex-row">
      <Button type="submit" disabled={$submitting}>
        Update Password
        {#if $submitting}
          <Spinner />
        {/if}
      </Button>
    </div>
  </form>
</div>
