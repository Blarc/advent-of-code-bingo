import {AsyncPipe, NgClass, NgForOf, NgIf} from '@angular/common';
import {Component, OnInit, inject} from '@angular/core';

import {Subject, tap} from 'rxjs';

import {BingoCardDto} from '../../core/api/model/bingoCardDto.model';
import {UserDto} from '../../core/api/model/userDto.model';
import {BingoApiService} from '../../core/api/service/bingo-api.service';
import {BingoGridComponent} from '../../core/components/bingo-grid/bingo-grid.component';
import {UserCountPipe} from '../../core/pipes/user-count.pipe';
import {AuthService} from '../../core/services/auth.service';

@Component({
    standalone: true,
    selector: 'app-bingo-page',
    templateUrl: 'bingo.page.component.html',
    imports: [NgForOf, NgIf, NgClass, AsyncPipe, UserCountPipe, BingoGridComponent]
})
export class BingoPageComponent implements OnInit {
    private authService = inject(AuthService);
    private bingoApiService = inject(BingoApiService);

    private user?: UserDto | null;
    public bingoCardsSubject = new Subject<BingoCardDto[]>();

    ngOnInit(): void {
        this.authService.userSubject$
            .pipe(
                tap(u => {
                    this.user = u;
                    this.fetchBingoCards();
                })
            )
            .subscribe();
    }

    private fetchBingoCards() {
        this.bingoApiService.getAllBingoCards().subscribe(cards => {
            if (this.user) {
                for (const card of cards) {
                    for (const userCard of this.user.bingo_cards) {
                        if (card.id === userCard.id) {
                            card.selected = true;
                        }
                    }
                }
            }
            this.bingoCardsSubject.next(cards);
        });
    }
}
