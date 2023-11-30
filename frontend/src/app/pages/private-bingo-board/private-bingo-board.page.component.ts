import {AsyncPipe, NgClass, NgForOf, NgIf} from '@angular/common';
import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';

import {Subject, tap} from 'rxjs';

import {BingoBoardDto} from '../../core/api/model/bingoBoardDto.model';
import {BingoCardDto} from '../../core/api/model/bingoCardDto.model';
import {UserDto} from '../../core/api/model/userDto.model';
import {BingoApiService} from '../../core/api/service/bingo-api.service';
import {BingoGridComponent} from '../../core/components/bingo-grid/bingo-grid.component';
import {UserCountPipe} from '../../core/pipes/user-count.pipe';
import {AuthService} from '../../core/services/auth.service';

@Component({
    standalone: true,
    selector: 'app-private-bingo-board',
    styleUrls: ['private-bingo-board.page.component.scss'],
    imports: [AsyncPipe, NgForOf, UserCountPipe, NgClass, NgIf, BingoGridComponent],
    templateUrl: 'private-bingo-board.page.component.html'
})
export class PrivateBingoBoardPageComponent implements OnInit {
    private user?: UserDto | null;
    private readonly boardUuid: string;
    public bingoBoardSubject = new Subject<BingoBoardDto>();
    public bingoCardsSubject = new Subject<BingoCardDto[]>();

    constructor(
        private route: ActivatedRoute,
        private authService: AuthService,
        private bingoApiService: BingoApiService
    ) {
        this.boardUuid = this.route.snapshot.params['uuid'];
    }

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
        this.bingoApiService.getBingoBoard(this.boardUuid).subscribe(board => {
            this.bingoBoardSubject.next(board);
            this.bingoCardsSubject.next(this.bingoCardsWithSelection(board.bingo_cards));
        });
    }

    private bingoCardsWithSelection(cards: BingoCardDto[]): BingoCardDto[] {
        return cards.map(card => {
            if (this.user?.bingo_cards && this.user?.bingo_cards.some(userCard => userCard.id === card.id)) {
                return {...card, selected: true};
            }
            return card;
        });
    }
}
