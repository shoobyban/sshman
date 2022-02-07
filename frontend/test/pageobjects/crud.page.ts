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

    public get btnListItem(): ((selector:String) => ChainablePromiseElement<Promise<WebdriverIO.Element>>) {
        return (selector:String) => $('#list-items '+selector);
    }

    public get modalAdd(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#add-modal');
    }

    public get modalEdit(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#edit-modal');
    }

    public get modalDelete(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#delete-modal');
    }

    public get searchField(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#search-input');
    }

    public get btnDeleteItem(): ChainablePromiseElement<Promise<WebdriverIO.Element>> {
        return $('#delete-item');
    }

    public async addItem (searchKey: string, item: object): Promise<void> {
        await this.btnAddItems.click();
        await expect(this.modalAdd).toBeDisplayedInViewport();
        for (let key in item) {
            let ifield = await this.inputField('add-'+key);
            let className = await ifield.getAttribute('class');
            let elementType = await ifield.getTagName();
            if (elementType == 'div' && className == 'multiselect' ) {
                await ifield.click();
                await ifield.keys(item[key]);
                await browser.pause(100);
                await ifield.keys("\uE007");
                await browser.pause(100); // doesn't work on single vueform/multiselect yet
            } else {
                ifield.setValue(item[key]);
            }
        }
        await this.btnAddSave.click();
        await browser.pause(100)
        await expect(this.modalAdd).not.toBeDisplayedInViewport();
        await this.searchField.setValue(item[searchKey]);
        await browser.pause(100);
        const items = $('#list-items');
        await expect(items).toHaveChildren({eq: 1});
    }

    public async editItem (searchValue: string, searchKey: string, item: object): Promise<void> {
        await this.searchField.setValue(searchValue);
        await browser.pause(100)
        const items = $('#list-items')
        await expect(items).toHaveChildren({eq: 1});
        await this.btnListItem('button.edit-item').click();
        await expect(this.modalEdit).toBeDisplayedInViewport();
        for (let key in item) {
            await this.inputField('edit-'+key).setValue(item[key]);
        }
        await this.btnEditSave.click();
        await browser.pause(100)
        await expect(this.modalEdit).not.toBeDisplayedInViewport();
        await this.searchField.setValue(item[searchKey]);
        await browser.pause(100);
        await expect(items).toHaveChildren({eq: 1});
    }

    public async deleteItem (searchValue: string): Promise<void> {
        await this.searchField.setValue(searchValue);
        await browser.pause(100)
        const items = $('#list-items')
        await expect(items).toHaveChildren({eq: 1});
        await this.btnListItem('button.delete-item').click();
        await expect(this.modalDelete).toBeDisplayedInViewport();
        await this.btnDeleteItem.click();
        await browser.pause(100);
        await expect(this.modalDelete).not.toBeDisplayedInViewport();
        await this.searchField.setValue(searchValue);
        await browser.pause(100);
        await expect(items).not.toHaveChildren({});
    }

    public async search (find: string): Promise<string[]> {
        await this.searchField.setValue(find);
        await browser.pause(1000);
        const items = $('#list-items');
        await expect(items).toHaveChildren({eq: 1});
        return await items.$$('tr').map(async (item) => {
            return await item.getAttribute('data-rowid');
        });
    }

}
