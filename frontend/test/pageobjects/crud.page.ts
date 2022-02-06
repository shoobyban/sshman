import { ChainablePromiseElement } from 'webdriverio';

import Page from './page';

export default class CrudPage extends Page {

    public get btnAddItems(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#add-items');
    }

    public get btnAddSave(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#add-save');
    }

    public get btnEditSave(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#edit-save');
    }

    public get inputField(): ((field:String) => ChainablePromiseElement<Promise<WebdriverIO.Element>>) {
        return (field:String) => $('#'+field);
    }

    public get btnEditItem(): ((selector:String) => ChainablePromiseElement<Promise<WebdriverIO.Element>>) {
        return (selector:String) => $('#list-items'+selector);
    }

    public get modalAdd(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#add-modal');
    }

    public get modalEdit(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#edit-modal');
    }

    public get searchField(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#search-input');
    }

    public async addItem (searchKey: string, item: object): Promise<void> {
        await this.btnAddItems.click();
        await expect(this.modalAdd).toBeDisplayedInViewport();
        for (let key in item) {
            await this.inputField('add-'+key).setValue(item[key]);
        }
        await this.btnAddSave.click();
        await browser.pause(100)
        await expect(this.modalAdd).not.toBeDisplayedInViewport();
        await this.searchField.setValue(item[searchKey]);
        await browser.pause(100);
        const items = $('#list-items');
        await expect(items).toHaveChildren(1);
    }

    public async editItem (searchValue: string, item: object): Promise<void> {
        await this.searchField.setValue(searchValue);
        await browser.pause(100)
        const items = $('#list-items')
        await expect(items).toHaveChildren(1);
        await this.btnEditItem(' :first-child button').click();
        await expect(this.modalEdit).toBeDisplayedInViewport();
        for (let key in item) {
            await this.inputField('edit-'+key).setValue(item[key]);
        }
        await this.btnEditSave.click();
        await browser.pause(100)
        await expect(this.modalEdit).not.toBeDisplayedInViewport();
    }

    public async search (find: string): Promise<string[]> {
        await this.searchField.setValue(find);
        await browser.pause(1000);
        const items = $('#list-items');
        await expect(items).toHaveChildren(1);
        return await items.$$('tr').map(async (item) => {
            return await item.getAttribute('data-rowid');
        });
    }

}
