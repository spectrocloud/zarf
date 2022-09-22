import { expect, test } from '@playwright/test';

const checkbox = 'input[type=checkbox]';

test.beforeEach(async ({ page }) => {
	page.on('pageerror', (err) => console.log(err.message));
});

test.describe('initialize a zarf cluster', () => {
	test('configure the init package', async ({ page }) => {
		await page.goto('/auth?token=insecure&next=/initialize/configure');

		// Stepper
		await expect(page.locator('.stepper :text("1 Configure") .step-icon')).toHaveClass(/primary/);
		await expect(page.locator('.stepper :text("2 Review") .step-icon')).toHaveClass(/disabled/);
		await expect(page.locator('.stepper :text("3 Deploy") .step-icon')).toHaveClass(/disabled/);

		// Package details
		await expect(page.locator('text=Package Type ZarfInitConfig')).toBeVisible();
		await expect(
			page.locator('text=Meta data Name: init Description: Used to establish a new Zarf cluster')
		).toBeVisible();

		// Components (check most functionaliy with k3s component)
		let k3s = page.locator('.accordion:has-text("k3s (Optional)")');
		await expect(k3s.locator(checkbox)).toBeEnabled();
		await expect(
			page.locator('.component-accordion-header:has-text("*** REQUIRES ROOT *** Install K3s")')
		).toBeVisible();
		await expect(k3s.locator('code')).toBeHidden();
		await k3s.locator('button').click();
		await expect(k3s.locator('code')).toBeVisible();
		await expect(k3s.locator('code:has-text("name: k3s")')).toBeVisible();

		// Check remaining components for deploy states
		await validateRequiredCheckboxes(page);

		let logging = page.locator('.accordion:has-text("logging (Optional)")');
		await expect(logging.locator(checkbox)).toBeEnabled();

		let gitServer = page.locator('.accordion:has-text("git-server (Optional)")');
		await expect(gitServer.locator(checkbox)).toBeEnabled();

		await page.locator('text=review deployment').click();
		await expect(page).toHaveURL('/initialize/review');
	});

	test('review the init package', async ({ page }) => {
		await page.goto('/auth?token=insecure&next=/initialize/review');

		await validateRequiredCheckboxes(page);
	});
});

async function validateRequiredCheckboxes(page) {
	// Check remaining components for deploy states
	let injector = page.locator('.accordion:has-text("zarf-injector (Required)")');
	await expect(injector.locator(checkbox)).toBeDisabled();
	await expect(injector.locator(checkbox)).toBeChecked();

	let seedRegistry = page.locator('.accordion:has-text("zarf-seed-registry (Required)")');
	await expect(seedRegistry.locator(checkbox)).toBeDisabled();
	await expect(seedRegistry.locator(checkbox)).toBeChecked();

	let registry = page.locator('.accordion:has-text("zarf-registry (Required)")');
	await expect(registry.locator(checkbox)).toBeDisabled();
	await expect(registry.locator(checkbox)).toBeChecked();
}
