import {NgClass, NgForOf, NgIf} from '@angular/common';
import {Component, OnInit} from '@angular/core';

import {BingoCardDto} from '../../core/api/model/bingoCardDto';
import {UserDto} from '../../core/api/model/userDto';
import {ApiService} from '../../core/api/service/api.service';

@Component({
    standalone: true,
    selector: 'app-bingo-page',
    templateUrl: 'bingo.page.component.html',
    imports: [NgForOf, NgIf, NgClass],
    styleUrls: ['bingo.page.component.scss']
})
export class BingoPageComponent implements OnInit {
    public user: UserDto | undefined;
    public bingoCards: BingoCardDto[] = [];

    constructor(private apiService: ApiService) {}

    ngOnInit(): void {
        this.fetchBingoCards();
    }

    private fetchBingoCards(): void {
        this.apiService.getAllBingoCards().subscribe({
            next: bingoCards => (this.bingoCards = bingoCards),
            error: e => console.log(e)
        });
    }

    clickBingoCard(id: number): void {
        this.apiService.clickBingoCard(id).subscribe({
            next: user => (this.user = user),
            error: e => console.log(e)
        });
    }

    userHasBingoCardClass(id: number): string {
        if (this.user && this.user.bingo_cards.some(bingoCard => bingoCard.id === id)) {
            return 'box-yellow';
        }
        return 'box';
    }
}
